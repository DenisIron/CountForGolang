package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	var (
		wait        sync.WaitGroup //счетчик
		totalResult = 0
		k           = 5
	)
	finalAllUrls := make(chan bool) // канал для сигнала о том, что все горутины выполнены
	chUr := make(chan string, k)    // буферизированный канал для сайтов

	go urlsInChar(chUr) //Заполняем канал сайтами

	for range chUr {
		wait.Add(1)                           // добавляет 1 к счётчику из пакета sync
		go countGo(&totalResult, chUr, &wait) //Запускаем горутины для 5 сайтов
	}

	go allUrlsFinish(finalAllUrls, &wait) //Функция для проверки выполнения всех горутин
	<-finalAllUrls                        //Подаем сигнал о том, что все горутины выполнены
	fmt.Printf("Total: %d\n", totalResult)
}

func allUrlsFinish(finalAllUrls chan bool, wait *sync.WaitGroup) {
	wait.Wait() // Ждем завершения всех горутин, когда счетчик равен 0
	finalAllUrls <- true
}
func countGo(totalResult *int, chUr chan string, wait *sync.WaitGroup) {
	url := <-chUr // Передаем в переменную 1 сайт
	resp, err := http.Get(url)
	er(err)
	site, err := ioutil.ReadAll(resp.Body)
	er(err)
	count := strings.Count(string(site), "Go") //Считает количество вхождений на сайте
	fmt.Printf("Count for %s = %d\n", url, count)
	*totalResult += count //Суммирует вхождения на всех заданных сайтах
	wait.Done()           //уменьшает счётчик на 1
}

func urlsInChar(chUr chan string) {
	urlsStdin := bufio.NewReader(os.Stdin)
	for {
		urlsEr, err := urlsStdin.ReadString('\n')
		if err == nil {
			urlsEr = strings.Replace(urlsEr, "\n", "", -1)
			chUr <- urlsEr
		} else {
			er(err)
		}
	}
}

func er(err error) {
	if err != nil {
		fmt.Println("Error")
	}
}
