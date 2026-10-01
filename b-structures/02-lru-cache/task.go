package lru

// Задача: LRU-кэш (Least Recently Used). LeetCode 146, самая частая
// «большая» задача Go-собесов.
//
// Кэш фиксированной ёмкости. При переполнении вытесняется элемент,
// к которому дольше всего не обращались. Обращение — это и Get, и Set.
//
// Требование: Get и Set — за O(1).
//
// Проверка: go test ./...

type LRU struct {
	// TODO
}

// NewLRU — конструктор кэша ёмкостью capacity элементов.
func NewLRU(capacity int) *LRU {
	// TODO
	return &LRU{}
}

// Get возвращает значение и признак наличия;
// найденный элемент становится самым свежим.
func (c *LRU) Get(key any) (any, bool) {
	// TODO
	return nil, false
}

// Set вставляет или обновляет элемент; он становится самым свежим.
// При переполнении вытесняется элемент, к которому дольше всего не обращались.
func (c *LRU) Set(key, value any) {
	// TODO
}

// Len — текущее число элементов.
func (c *LRU) Len() int {
	// TODO
	return 0
}
