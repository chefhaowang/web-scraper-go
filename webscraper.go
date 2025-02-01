package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
	pb "web-scraper-go/scraperpb"
)

func ScrapeTopNews() []*pb.NewsArticle {
	chromeDriverPort := 52777
	caps := selenium.Capabilities{"browserName": "chrome"}

	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", chromeDriverPort))
	if err != nil {
		log.Fatalf("Error starting WebDriver: %v", err)
	}
	defer driver.Quit()

	url := "https://www.theverge.com/tech"
	if err := driver.Get(url); err != nil {
		log.Fatalf("Failed to load page: %v", err)
	}

	time.Sleep(5 * time.Second)

	titles, _ := driver.FindElements(selenium.ByCSSSelector, "a._1lkmsmo1")
	authors, _ := driver.FindElements(selenium.ByCSSSelector, "span._1lldluw2 span")
	dates, _ := driver.FindElements(selenium.ByCSSSelector, "span.duet--article--timestamp time")
	urls, _ := driver.FindElements(selenium.ByCSSSelector, "a._1lkmsmo1")

	var articles []*pb.NewsArticle

	for i := 0; i < 3; i++ {
		title, _ := titles[i].Text()
		author, _ := authors[i].Text()
		date, _ := dates[i].Text()
		link, _ := urls[i].GetAttribute("href")

		articles = append(articles, &pb.NewsArticle{
			Title:  title,
			Author: author,
			Date:   date,
			Url:    link,
		})
	}

	return articles
}