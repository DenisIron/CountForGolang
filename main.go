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
		wait sync.WaitGroup
		k    = 5
	)
	chUr := make(chan string, k)
	chEr := make(chan string)
	chCount := make(chan int)
	go urlsInChar(chUr)

	for i := 0; i < k; i++ {
		wait.Add(1)
		go func() {
			for url := range chUr {
				count := countGoOnSite(getSiteBody(url, chEr)) //Считает количество вхождений "Go" на сайте
				chCount <- count                               //Передаем количество вхождений в канал, для дальнейшего подсчета
				go printCount(url, count)
			}
			wait.Done()
		}()
	}
	go errHandling(chEr)
	wait.Wait()
	result := allCount(chCount) //Функция для подсчета общего кол-ва всех вхождений на сайтах
	printResult(result)
}

func getSiteBody(url string, chEr chan string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		chEr <- "Error for URL: " + string(url)
	}
	site, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		chEr <- "Error in Body of URL: " + string(url)
	}
	return site
}

func countGoOnSite(site []byte) int {
	count := strings.Count(string(site), "Go")
	return count
}

func printCount(url string, count int) {
	fmt.Printf("Count for %s = %d\n", url, count)
}

func urlsInChar(chUr chan string) { // +добавить передачу ошибки
	urlsStdin := bufio.NewReader(os.Stdin)
	for {
		urlsEr, err := urlsStdin.ReadString('\n')
		if err != nil {
			break
		}
		urlsEr = strings.Replace(urlsEr, "\n", "", -1)
		chUr <- urlsEr //передаем в буферизированный канал первые 5 сайтов
	}
}

func allCount(chCount chan int) int {
	totalResult := 0
	for range chCount { //Исправить!!!!
		totalResult += <-chCount
	}
	return totalResult
}

func printResult(result int) {
	fmt.Printf("Total: %d\n", result)
}

func errHandling(chEr chan string) {
	url string
}

/*Возможные варианты функций обработки ошибки:
func fatalErr(err error) {
	log.Fatal("Aborting: ", err)
}
func er(err error) {
	if err != nil {
		panic(err)
	}
}

func errors(err os.Error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
*/
