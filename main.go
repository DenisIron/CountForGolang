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
		url         string //строка для хранения 1 сайта
	)
	finalOne := make(chan bool)  // канал для сигнала о том, что 1 горутина выполнена
	chUr := make(chan string, k) // буферизированный канал для 5 сайтов
	chCount := make(chan int)    // канал для количества вхождений на каждом сайте

	go urlsInChar(chUr) // Заполняем канал сайтами

	for i := 0; i < k; i++ { //ограничиваем количество гороутин
		// с целью избавления от ошибки связанной с sync добавляем ф-цию в main
		go func() {
			for range chUr {
				wait.Add(1)  // добавляет 1 к счётчику из пакета sync
				url = <-chUr // передаем из буферизированного канала 1 сайт в переменную
				site := getBodySite(url)
				count := countGoOnSite(site) //Считает количество вхождений на сайте
				chCount <- count             //Передаем количество вхождений в канал, для дальнейшего подсчета
				printCount(url, count)       //Печатаем результат в консоли
				wait.Done()                  //Уменьшает счётчик на 1
				finalOne <- true
			}
		}()
	}
	wait.Wait()                              // Ждем завершения всех горутин, когда счетчик равен 0
	result := allCount(chCount, totalResult) //Функция для подсчета общего кол-ва всех вхождений на сайтах
	fmt.Printf("Total: %d\n", result)
}

func getBodySite(url string) []byte {
	resp, err := http.Get(url)
	er(err)
	site, err := ioutil.ReadAll(resp.Body)
	er(err)
	return site
}

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
			chUr <- urlsEr //передаем в буферизированный канал первые 5 сайтов
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
