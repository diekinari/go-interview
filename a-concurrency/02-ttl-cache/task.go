package ttlcache

import (
	"context"
	"time"
)

// Задача: кэш с TTL.
//
// Каждая запись живёт ttl с момента Set; повторный Set обновляет срок.
//
// Требования к поведению:
//   - Get по живой записи возвращает её значение;
//   - Get по протухшей записи возвращает (nil, false) — протухшее
//     не отдаётся никогда;
//   - протухшая запись физически исчезает не позже чем через ttl после
//     истечения, ДАЖЕ ЕСЛИ к ней больше никто не обращается (тесты
//     проверяют это через Len);
//   - кэш безопасен из множества горутин.
//
// Проверка: go test -race ./...

type TTLCache struct {
	// TODO
}

// NewTTLCache — конструктор; ctx ограничивает время жизни кэша.
func NewTTLCache(ctx context.Context, ttl time.Duration) *TTLCache {
	// TODO
	return &TTLCache{}
}

// Get возвращает значение и признак наличия; протухшая запись
// считается отсутствующей.
func (c *TTLCache) Get(key string) (any, bool) {
	// TODO
	return nil, false
}

// Set записывает значение; срок жизни — ttl с этого момента.
func (c *TTLCache) Set(key string, value any) {
	// TODO
}

// Len — сколько записей физически хранится в данный момент
// (нужен тестам для проверки требования про накопление).
func (c *TTLCache) Len() int {
	// TODO
	return 0
}
