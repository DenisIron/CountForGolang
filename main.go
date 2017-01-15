package main

import (
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
	var arrayUrls []string //Создаем slice, заполняемый вводимыми сайтами
	urls(arrayUrls)        //Заполняем slice сайтами

	for i := range arrayUrls {
		url := arrayUrls[i]
		if arrayUrls[i] == "" {
			break
		}
		if gorout < k {
			wait.Add(1)
			go countGo(url, &totalResult, &gorout) //Запускаем горутины для 5 сайтов
		} else {
			wait.Wait() //Если горутин больше 5, то ждем пока хотя бы одна из них не выполнится
			wait.Add(1)
			go countGo(url, &totalResult, &gorout)
		}

	}

	fmt.Printf("Total: %d\n", totalResult)
}

func countGo(url string, totalResult *int, gorout *int) {
	*gorout++
	resp, err := http.Get(url)
	site, err := ioutil.ReadAll(resp.Body)
	er(err)
	count := strings.Count(string(site), "Go") //Считает количество вхождений на сайте
	fmt.Printf("Count for %s = %d\n", url, count)
	*totalResult += count //Суммируем вхождения на всех заданных сайтах
	*gorout--
}

func urls(urls []string) {
	//scanner := bufio.NewScanner(os.Stdin)    //Считываем строку с сайтами из командной строки
	urlsStdin, err := ioutil.ReadAll(os.Stdin) //Проверяем содержимое командной строки
	if err != nil {
		urls = strings.Split(string(urlsStdin), "\n")
	} else {
		er(err)
	}
}
func er(err error) {
	if err != nil {
		panic(er)
	}
}
