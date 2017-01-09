package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

func main() {
	var (
		totalResult = 0
		k           = 5
		gorout      = 0
		wait        sync.WaitGroup
	)
	searchGo := "Go"             //Прописываем какие вхождения будем искать на сайтах
	urls := make([]string, 1000) //Создаем slice, заполняемый вводимыми сайтами
	for i := range urls {
		fmt.Scanln(&urls[i]) //Вводим URL, в котором будем искать количество вхождений
		url := urls[i]
		if gorout < k {
			wait.Add(1)
			go countGo(url, &totalResult, &gorout, searchGo) //Запускаем горутины для 5 сайтов
		} else {
			wait.Wait() //Если горутин больше 5, то ждем пока хотя бы одна из них не выполнится
			wait.Add(1)
			go countGo(url, &totalResult, &gorout, searchGo)
		}
		if urls[i] == "" {
			break
		}
	}
	fmt.Printf("Total: %d\n", totalResult)
}

func countGo(url string, totalResult *int, gorout *int, searchGo string) {
	*gorout++
	resp, err := http.Get(url)
	site, err := ioutil.ReadAll(resp.Body)
	er(err)
	count := strings.Count(string(site), searchGo) //Считает количество вхождений на сайте
	fmt.Printf("Count for %s = %d\n", url, count)
	*totalResult += count //Суммируем вхождения на всех заданных сайтах
	*gorout--
}

func er(err error) {
	if err != nil {
		panic(er)
	}
}
