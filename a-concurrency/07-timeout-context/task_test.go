package timeoutctx

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"
)

func TestFastResult(t *testing.T) {
	res, err := FetchData(context.Background(), func() string {
		return "data"
	})
	if err != nil || res != "data" {
		t.Fatalf("ожидали (data, nil), получили (%q, %v)", res, err)
	}
}

func TestTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	res, err := FetchData(ctx, func() string {
		time.Sleep(500 * time.Millisecond)
		return "late"
	})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatalf("ожидали ошибку по таймауту, получили (%q, nil)", res)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("ожидали context.DeadlineExceeded, получили %v", err)
	}
	if elapsed > 200*time.Millisecond {
		t.Fatalf("возврат занял %v — FetchData дождался slow() вместо отмены", elapsed)
	}
}

func TestNoGoroutineLeak(t *testing.T) {
	before := runtime.NumGoroutine()

	for i := 0; i < 20; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		FetchData(ctx, func() string {
			time.Sleep(50 * time.Millisecond)
			return "late"
		})
		cancel()
	}

	// Даём фоновым горутинам время завершиться.
	time.Sleep(300 * time.Millisecond)
	after := runtime.NumGoroutine()

	if after > before+2 {
		t.Fatalf("утечка горутин: было %d, стало %d — фоновые горутины не завершились", before, after)
	}
}
