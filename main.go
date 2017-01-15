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
		gorout      = 0
		wait        sync.WaitGroup
	)
	finalAllUrls := make(chan bool) // канал для сигнала о том, что все горутины выполнены
	finalOneUrls := make(chan bool) // канал для сигнала о выполнении 1 горутины
	var sliceUrls []string          //Создаем slice, заполняемый вводимыми сайтами

	urls(sliceUrls)            //Заполняем slice сайтами
	lenSlice := len(sliceUrls) //Количество всех сайтов
	for i := range sliceUrls {
		url := sliceUrls[i] //Заносим в отдельную переменную 1 сайт
		if gorout < k {
			wait.Add(1)
			go countGo(url, &totalResult, &gorout, finalOneUrls) //Запускаем горутины для 5 сайтов
		} else {
			wait.Wait() //Если горутин больше 5, то ждем пока хотя бы одна из них не выполнится
			wait.Add(1)
			go countGo(url, &totalResult, &gorout, finalOneUrls)
		}
	}
	go allUrls(finalOneUrls, finalAllUrls, lenSlice) //Функция для проверки выполнения всех горутин
	<-finalAllUrls                                   // Подаем сигнал о том, что все горутины выполнены
	fmt.Printf("Total: %d\n", totalResult)
}

func allUrls(finalOneUrls chan bool, finalAllUrls chan bool, lenSlice int) {
	for i := 0; i < lenSlice; i++ {
		<-finalOneUrls
	}
	finalAllUrls <- true
}
func countGo(url string, totalResult *int, gorout *int, oneUr chan bool) {
	*gorout++
	resp, err := http.Get(url)
	site, err := ioutil.ReadAll(resp.Body)
	er(err)
	count := strings.Count(string(site), "Go") //Считает количество вхождений на сайте
	fmt.Printf("Count for %s = %d\n", url, count)
	*totalResult += count //Суммируем вхождения на всех заданных сайтах
	*gorout--
	oneUr <- true
}

func urls(sliceUrls []string) []string {
	urlsStdin := bufio.NewReader(os.Stdin) //Считываем строку с сайтами из командной строки
	//urlsStdin, err := ioutil.ReadAll(os.Stdin) //Второй возможный вариант
	urlsEr, err := urlsStdin.ReadString('\n') //Проверяем содержимое командной строки
	if err == nil {
		sliceUrls = strings.Split(urlsEr, "\n") //Полученную на входе строку преобразуем в slice с Urls
	}
	return sliceUrls
}

func er(err error) {
	if err != nil {
		panic(er)
	}
}
