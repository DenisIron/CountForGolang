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
		totalResult int
		k           = 5
		//urlsSlice   = []string{} //slice для всех Urls
		url string //строка для хранения 1 сайта
	)
	//finalAllUrls := make(chan bool) // канал для сигнала о том, что все горутины выполнены
	finalOne := make(chan bool)  // канал для сигнала о том, что 1 горутина выполнена
	chUr := make(chan string, k) // буферизированный канал для 5 сайтов
	chCount := make(chan int)    // канал для количества вхождений на каждом сайте

	go urlsInChar(chUr) // Заполняем канал сайтами

	/*	for range urlsSlice {
		urlsSlice[i] = <-chUr // Передаем в переменную 1 сайт
	}*/

	for {
		url = <-chUr // !будет ошибка из-за того, что канал пока пуст
		wait.Add(1)  // добавляет 1 к счётчику из пакета sync

		// с целью избавления от ошибки связанной с sync, добавляем ф-цию в main
		go func(url string, chCount chan int, finalOne chan bool) { //Запускаем горутины для 5 сайтов
			defer wait.Done() //Уменьшает счётчик на 1
			resp, err := http.Get(url)
			er(err)
			site, err := ioutil.ReadAll(resp.Body)
			er(err)
			count := countGoOnSite(site) //Считает количество вхождений на сайте
			chCount <- count             //Передаем количество вхождений в канал, для дальнейшего подсчета
			printCount(url, count)       //Печатаем результат в консоли
			finalOne <- true
		}(url, chCount, finalOne)
		wait.Wait() // Ждем завершения всех горутин, когда счетчик равен 0
		break       //После завершения всех гороутин выходим из цикла
	}

	/*for _, url := range urlsSlice {
		wait.Add(1)                    // добавляет 1 к счётчику из пакета sync
		go countGo(url, wait, chCount) //Запускаем горутины для 5 сайтов
	}*/

	//go allUrlsFinish(finalAllUrls, &wait) //Функция для проверки выполнения всех горутин
	//<-finalAllUrls                        //Подаем сигнал о том, что все горутины выполнены
	//wait.Wait() // Ждем завершения всех горутин, когда счетчик равен 0

	allCount(chCount, totalResult) //Функция для подсчета общего кол-ва всех вхождений на сайтах
	fmt.Printf("Total: %d\n", totalResult)
}

/*func allUrlsFinish(finalAllUrls chan bool, wait *sync.WaitGroup) {
	wait.Wait() // Ждем завершения всех горутин, когда счетчик равен 0
	finalAllUrls <- true
}*/

func countGoOnSite(site []byte) int {
	count := strings.Count(string(site), "Go")
	return count
}

func printCount(url string, count int) {
	fmt.Printf("Count for %s = %d\n", url, count)
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

func allCount(chCount chan int, totalResult int) int {
	for range chCount {
		totalResult += <-chCount
	}
	return totalResult
}

func er(err error) {
	if err != nil {
		panic(err)
	}
}
