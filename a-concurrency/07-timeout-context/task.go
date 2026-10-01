package timeoutctx

import "context"

// Задача: операция с таймаутом.
//
// slow — медленная операция, которую нельзя прервать (ctx она не принимает).
// FetchData вызывает её с учётом ctx.
//
// Требования:
//   - slow() успела до отмены ctx -> вернуть её результат и nil;
//   - ctx отменён/истёк раньше -> вернуть ("", ctx.Err()) сразу,
//     не дожидаясь завершения slow();
//   - FetchData не оставляет после себя ничего, что живёт дольше
//     самой slow() (тест считает runtime.NumGoroutine).
//
// Проверка: go test -race ./...


func FetchData(ctx context.Context, slow func() string) (string, error) {
	// TODO
	return "", nil
}
