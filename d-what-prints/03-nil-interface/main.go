package main

import "fmt"

// Что выведет программа? (реальный собес Lamoda)
// Три строки печати — какие из них выполнятся?
//
// Сначала ответь, потом скролль к разбору.

type Seller interface {
	GetID() int64
}

type Lamoda struct{}

func (l Lamoda) GetID() int64 { return 1 }

func main() {
	var seller Seller
	if seller == nil {
		fmt.Println("nil interface")
	}

	var lamoda *Lamoda
	if lamoda == nil {
		fmt.Println("nil struct")
	}

	seller = lamoda
	if seller == nil {
		fmt.Println("nil assignment")
	}
}

//
//
//
//
//
//
//
//
//
//
// ================================ РАЗБОР ================================
//
// Вывод:
//   nil interface
//   nil struct
//
// Третья строка НЕ печатается.
//
// Интерфейсное значение — пара (динамический тип, значение) и равно nil,
// только когда nil ОБЕ части.
//   var seller Seller        -> (nil, nil)      -> == nil, печать
//   var lamoda *Lamoda       -> nil-указатель   -> == nil, печать
//   seller = lamoda          -> (*Lamoda, nil)  -> тип НЕ nil -> != nil
//
// Это «типизированный nil» — источник классического бага: функция
// с возвращаемым типом error возвращает nil-указатель конкретного
// типа ошибки, и у вызывающего if err != nil срабатывает на «пустой»
// ошибке. Правило: нет ошибки — возвращай литеральный nil.
//
// Вопрос вдогонку: удовлетворяет ли *Lamoda интерфейсу Seller,
// если метод объявлен на value receiver? (Да: набор методов *T
// включает методы T.)
