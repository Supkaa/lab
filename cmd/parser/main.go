package main

import (
	"database/sql"
	"fmt"
	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
)

func main() {
	db, err := sql.Open("sqlite3", "./storage/lab2.db")
	if err != nil {
		log.Panicf("fail to connect db: %s", err.Error())
	}

	pageCount := 20
	rawNews := make(chan []string, pageCount)
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(pageCount)

	go parseNew(rawNews, db, done)

	for page := 1; page <= pageCount; page++ {
		go func(page int) {
			defer wg.Done()
			rawNews <- parseNewsPage(page)
		}(page)
	}

	wg.Wait()
	close(rawNews)

	<-done

	fmt.Println("Главная горутина завершила выполнение.")
}

func parseNewsPage(page int) []string {
	var news []string
	client := resty.New().GetClient()

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://news.tyuiu.ru/sections/novosti-tiu?page=%d", page), nil)
	response, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		panic(err)
	}

	n := doc.Find(".grid")
	n.Children().Each(func(i int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if ok {
			news = append(news, href)
		}
	})

	return news
}

func parseNew(in <-chan []string, db *sql.DB, done chan<- struct{}) {
	baseId := uuid.MustParse("1d03b612-d864-49c3-90bd-0ea2f5114d07")

	for urls := range in {
		client := resty.New().GetClient()
		wg := sync.WaitGroup{}
		wg.Add(len(urls))

		for _, url := range urls {
			go func(url string) {
				defer wg.Done()
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				response, err := client.Do(req)

				if err != nil {
					panic(err)
				}

				defer response.Body.Close()
				if response.StatusCode != http.StatusOK {
					log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
				}

				doc, err := goquery.NewDocumentFromReader(response.Body)

				if err != nil {
					panic(err)
				}

				if _, err := db.Exec("INSERT INTO news(id, title, summary, image, created_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET title = $2, summary = $3, image = $4, created_at = $5", uuid.NewSHA1(baseId, []byte(url)), getTitle(doc), getDescription(doc), getImage(doc), getPublishedAt(doc)); err != nil {
					log.Fatal(err)
				}

				log.Printf("title: %s\n description: %s\n image: %s\n time: %s", getTitle(doc), getDescription(doc), getImage(doc), getPublishedAt(doc))
			}(url)
		}

		wg.Wait()
	}

	done <- struct{}{}
}

func getPublishedAt(doc *goquery.Document) time.Time {
	meta := doc.Find(".flex.items-center.gap-2").Eq(5)

	t, err := time.Parse("02.01.2006", strings.TrimSpace(meta.Text()))

	if err != nil {
		panic(err)
	}

	return t
}

func getTitle(doc *goquery.Document) string {
	meta := doc.Find(".font-extrabold.text-xl")

	return strings.TrimSpace(meta.Text())
}

func getDescription(doc *goquery.Document) string {
	meta := doc.Find("#content")

	return strings.TrimSpace(meta.Children().Eq(0).Text())
}

func getImage(doc *goquery.Document) string {
	meta := doc.Find("main")

	imgPath, ok := meta.Children().First().Children().Attr("src")

	if !ok {
		return ""
	}

	return strings.TrimSpace(imgPath)
}
