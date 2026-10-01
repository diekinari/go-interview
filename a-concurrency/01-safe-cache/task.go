package safecache

// Задача: потокобезопасный кэш.
//
// Реализуй кэш, безопасный для использования из множества горутин.
// Чтений ожидается сильно больше, чем записей.
//
// Проверка: go test -race ./...

type Cache struct {
	// TODO
}

// NewCache — конструктор.
func NewCache() *Cache {
	// TODO
	return &Cache{}
}

// Get возвращает значение и признак наличия ключа.
func (c *Cache) Get(key string) (any, bool) {
	// TODO
	return nil, false
}

// Set записывает значение по ключу (перезаписывает существующее).
func (c *Cache) Set(key string, value any) {
	// TODO
}

// Delete удаляет ключ; удаление отсутствующего ключа — не ошибка.
func (c *Cache) Delete(key string) {
	// TODO
}
