package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/badoux/goscraper"
	"mvdan.cc/xurls/v2"
)

// var (
// 	indexRe = regexp.MustCompile(`^\d+[.)]\s*`)
// )

func main() {
	err := run()
	if err != nil {
		log.Fatalf("** %+v", err)
	}
}

func run() error {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")

	urlsRe := xurls.Relaxed()

	today := time.Now().Format("02.01.2006")

	output, err := os.Create("output.txt")
	if err != nil {
		return err
	}
	defer output.Close()

	errors, err := os.Create("errors.txt")
	if err != nil {
		return err
	}
	defer errors.Close()

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// line = indexRe.ReplaceAllLiteralString(line, "")

		urls := urlsRe.FindAllString(line, -1)
		if len(urls) == 0 {
			fmt.Fprintf(errors, "No URLs found in line: %v\n", line)
		}
		for _, url := range urls {
			if !strings.Contains(url, "://") {
				url = "http://" + url
			}

			var title string

			s, err := goscraper.Scrape(url, 5)
			if err != nil {
				fmt.Fprintf(errors, "Cannot download %v : %v\n", url, err)
				title = "ОШИБКА ПРИ ПОЛУЧЕНИИ ЗАГОЛОВКА"
			} else {
				title = strings.TrimSpace(s.Preview.Title)
			}

			fmt.Fprintf(output, "%s [Электронный ресурс]. – Режим доступа: %s. – Заглавие с экрана. – (Дата обращения: %s).\n", title, url, today)
		}

	}

	return nil
}
