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
		totalResult = 0 //Итоговое кол-во вхождений
		k           = 5 //Макс. число горутин
		//gorout      = 0
		wait sync.WaitGroup
	)
	finalAllUrls := make(chan bool) // канал для сигнала о том, что все горутины выполнены
	finalOneUr := make(chan bool)   // канал для сигнала о выполнении 1 горутины
	chUr := make(chan string)       // канал для сайтов

	urls(chUr) //Заполняем канал сайтами

	for gorout := 0; gorout < k; gorout++ { // создаем замыкание, тем самым ограничив кол-во горутин
		go countGo(&totalResult, finalOneUr, chUr, &wait) //Запускаем горутины для 5 сайтов
		//<-finalOneUr                                      //Если горутин больше 5, то ждем пока хотя бы одна из них выполнится
	}
	go allUrls(finalAllUrls, &wait) //Функция для проверки выполнения всех горутин
	<-finalAllUrls                  // Подаем сигнал о том, что все горутины выполнены
	fmt.Printf("Total: %d\n", totalResult)
}

func allUrls(finalAllUrls chan bool, wait *sync.WaitGroup) {
	wait.Wait() // Ждем завершения всех горутин
	finalAllUrls <- true
}
func countGo(totalResult *int, oneUr chan bool, chUr chan string, wait *sync.WaitGroup) {
	wait.Add(1)
	url := <-chUr // Передаем в переменную 1 сайт
	resp, err := http.Get(url)
	er(err)
	site, err := ioutil.ReadAll(resp.Body)
	er(err)
	count := strings.Count(string(site), "Go") //Считает количество вхождений на сайте
	fmt.Printf("Count for %s = %d\n", url, count)
	*totalResult += count //Суммируем вхождения на всех заданных сайтах
	wait.Add(-1)
	oneUr <- true
}

func urls(chUr chan string) {
	urlsStdin := bufio.NewReader(os.Stdin) //Считываем строку с сайтами из командной строки
	for {
		urlsEr, err := urlsStdin.ReadString('\n') //Проверяем содержимое командной строки
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
		panic(er)
	}
}
