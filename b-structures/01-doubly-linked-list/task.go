package dlist

// Задача: двусвязный список с нуля (без container/list).
//
//
// Проверка: go test ./...

type Node struct {
	Value any
	Prev  *Node
	Next  *Node
}

type List struct {
	// TODO
}

// NewList — конструктор пустого списка.
func NewList() *List {
	// TODO
	return &List{}
}

// PushFront вставляет значение в начало и возвращает созданный узел.
func (l *List) PushFront(v any) *Node {
	// TODO
	return nil
}

// Remove удаляет узел из списка.
func (l *List) Remove(n *Node) {
	// TODO
}

// MoveToFront перемещает существующий узел в начало.
func (l *List) MoveToFront(n *Node) {
	// TODO
}

// Back возвращает последний узел (nil, если список пуст).
func (l *List) Back() *Node {
	// TODO
	return nil
}

// Values возвращает все значения от головы к хвосту (нужен тестам).
func (l *List) Values() []any {
	// TODO
	return nil
}
