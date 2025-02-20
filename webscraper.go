package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
	pb "web-scraper-go/scraperpb"

	"github.com/tebeka/selenium"
)

func cleanupChrome() {
    exec.Command("pkill", "-f", "chrome").Run()
    exec.Command("pkill", "-f", "chromedriver").Run()
    exec.Command("rm", "-rf", "/tmp/chrome-user-data-*").Run()
}



func ScrapeTopNews() []*pb.NewsArticle {
	chromeDriverPort := 52777
	caps := selenium.Capabilities{"browserName": "chrome"}

	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://3.104.120.8:%d", chromeDriverPort))
	if err != nil {
		log.Fatalf("Error starting WebDriver: %v", err)
	}

	// Call cleanup after each scraping session
	defer cleanupChrome()
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
		href, _ :=  urls[i].GetAttribute("href")


		// Combine base URL with href if needed
		link := href
		if !strings.HasPrefix(href, "http") {
			link = "https://www.theverge.com" + href
		}

		articles = append(articles, &pb.NewsArticle{
			Title:  title,
			Author: author,
			Date:   date,
			Url:    link,
		})
	}

	return articles
}

