package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	totalResult := 0
	gorout := 0
	startNext := make(chan bool)
	urls := make([]string, 1000) //Создаем срез, заполняемый вводимыми сайтами
	for i := range urls {
		fmt.Scanln(&urls[i]) //Вводим URL, в котором будем искать количество вхождений
		if urls[i] == "" {
			break
		}
		if gorout < 5 {
			gorout++
			go countGo(urls, i, &totalResult, startNext, &gorout) //Запускаем горутины для 5 сайтов
		} else {
			func() {
				<-startNext //Если горутин больше 5, то ждем пока хотя бы одна из них не выполнится
			}()
		}

	}
	defer fmt.Printf("Total: %d\n", totalResult)
}

func countGo(urls []string, i int, totalResult *int, startNext chan bool, gorout *int) {
	searchGo := "Go" // Прописываем какие вхождения будем искать на сайтах
	url := urls[i]
	resp, err := http.Get(url)
	er(err)
	site, err := ioutil.ReadAll(resp.Body)
	er(err)
	text := string(site)
	count := strings.Count(text, searchGo) //Считает количество вхождений на сайте
	fmt.Printf("Count for %s = %d\n", url, count)
	*totalResult += count //Суммируем вхождения на всех заданных сайтах
	*gorout--
	startNext <- true
}

func er(err error) {
	if err != nil {
		return
	}
}
