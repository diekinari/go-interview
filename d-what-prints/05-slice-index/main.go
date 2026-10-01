package main

import (
	"fmt"
	"sync"
)

// Что выведет программа? (реальный собес)
// Дойдёт ли выполнение до "Program completed."?
//
// Сначала ответь, потом скролль к разбору.

type listSKU []string

func (l listSKU) getLastSKU() string {
	return l[len(l)]
}

func main() {
	items := listSKU{
		"MP990099991",
		"MP900000002",
		"MP000000003",
		"MP000000004",
		"MP000000005",
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		lastItem := items.getLastSKU()
		fmt.Printf("Last SKU: %s\n", lastItem)
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Program completed.")
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
// Программа упадёт: "panic: runtime error: index out of range [5] with
// length 5". Ошибка в getLastSKU: последний элемент — l[len(l)-1],
// а l[len(l)] — выход за границу.
//
// Два усиливающих момента, которые надо проговорить:
//
// 1. Паника происходит В ГОРУТИНЕ. Непойманная panic в любой горутине
//    роняет ВЕСЬ процесс — main не «переживёт» чужую панику.
//    До "Program completed." выполнение не дойдёт.
//
// 2. wg.Done() стоит ПОСЛЕ кода, который может паниковать, и не через
//    defer. Даже если бы панику где-то recover-или, Done не выполнился бы
//    и wg.Wait() завис. Правильно: defer wg.Done() первой строкой горутины.
//
// Вопрос вдогонку: где можно поставить recover, чтобы программа
// не упала? (defer с recover ВНУТРИ той же горутины — recover из main
// чужую панику не ловит.)
