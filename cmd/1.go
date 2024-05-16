package main

import (
	"fmt"
	"sync"
)

func main() {
	channel := make(chan int)
	var wg sync.WaitGroup

	// Запускаем горутины пишущие в канал
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go writeToChannel(channel, i, &wg)
	}

	// Горутина читающая из канала
	go readFromChannel(channel)

	// Ждем завершения всех горутин пишущих в канал
	wg.Wait()

	// Закрываем канал, чтобы прекратить чтение
	close(channel)

	fmt.Println("Главная горутина завершила выполнение.")
}

func writeToChannel(ch chan<- int, num int, wg *sync.WaitGroup) {
	defer wg.Done()
	ch <- num
	fmt.Printf("Горутина записала в канал число: %d\n", num)
}

func readFromChannel(ch <-chan int) {
	for value := range ch {
		fmt.Printf("Прочитано из канала: %d\n", value)
	}
}
