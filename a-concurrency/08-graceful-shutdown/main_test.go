package main

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

const addr = "127.0.0.1:8080"

// bin — путь к собранной программе; собирается один раз в TestMain.
var bin string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "graceful")
	if err != nil {
		panic(err)
	}
	bin = filepath.Join(dir, "server")
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		os.Stderr.Write(out)
		os.RemoveAll(dir)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// server — запущенная программа и её вывод.
type server struct {
	cmd    *exec.Cmd
	out    bytes.Buffer
	exited chan error // сюда приходит результат cmd.Wait
}

// start запускает программу и ждёт, пока она начнёт принимать соединения.
func start(t *testing.T) *server {
	t.Helper()
	if conn, err := net.Dial("tcp", addr); err == nil {
		conn.Close()
		t.Fatalf("порт %s уже занят — освободи его перед тестом", addr)
	}

	s := &server{exited: make(chan error, 1)}
	s.cmd = exec.Command(bin)
	s.cmd.Stdout = &s.out
	s.cmd.Stderr = &s.out
	if err := s.cmd.Start(); err != nil {
		t.Fatal(err)
	}
	go func() { s.exited <- s.cmd.Wait() }()
	t.Cleanup(func() {
		s.cmd.Process.Kill()
		<-s.exited
		if t.Failed() {
			t.Logf("вывод программы:\n%s", s.out.String())
		}
	})

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if conn, err := net.Dial("tcp", addr); err == nil {
			conn.Close()
			return s
		}
		select {
		case err := <-s.exited:
			s.exited <- err // вернуть для Cleanup
			t.Fatalf("программа завершилась, не начав слушать %s: %v", addr, err)
		case <-time.After(50 * time.Millisecond):
		}
	}
	t.Fatalf("за 5 секунд программа не начала слушать %s", addr)
	return nil
}

// waitExit ждёт завершения программы не дольше within и проверяет код 0.
func (s *server) waitExit(t *testing.T, within time.Duration) {
	t.Helper()
	select {
	case err := <-s.exited:
		s.exited <- err // вернуть для Cleanup
		if err != nil {
			t.Fatalf("программа должна завершиться с кодом 0, а завершилась: %v", err)
		}
	case <-time.After(within):
		t.Fatalf("программа не завершилась за %v после сигнала", within)
	}
}

func TestIdleExitsImmediately(t *testing.T) {
	for _, sig := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run(sig.String(), func(t *testing.T) {
			s := start(t)
			if err := s.cmd.Process.Signal(sig); err != nil {
				t.Fatal(err)
			}
			s.waitExit(t, time.Second)
		})
	}
}

func TestWaitsForActiveRequest(t *testing.T) {
	s := start(t)

	type result struct {
		status int
		body   string
		err    error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := http.Get("http://" + addr + "/")
		if err != nil {
			done <- result{err: err}
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		done <- result{status: resp.StatusCode, body: string(body), err: err}
	}()

	time.Sleep(300 * time.Millisecond) // запрос уже в обработке
	if err := s.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}

	// Новые соединения после сигнала не принимаются.
	time.Sleep(200 * time.Millisecond)
	if conn, err := net.Dial("tcp", addr); err == nil {
		conn.Close()
		t.Fatal("после сигнала сервер всё ещё принимает новые соединения")
	}

	// Начатый запрос получает нормальный ответ.
	var r result
	select {
	case r = <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("начатый запрос не получил ответа")
	}
	if r.err != nil {
		t.Fatalf("начатый запрос оборвался: %v", r.err)
	}
	if r.status != http.StatusOK || r.body != "ok\n" {
		t.Fatalf("ожидали 200 и \"ok\\n\", получили %d и %q", r.status, r.body)
	}

	// Программа выходит сразу после того, как запрос обслужен.
	s.waitExit(t, time.Second)
}
