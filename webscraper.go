package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

func main() {
	chromeDriverPort := 52777

	// Start a Selenium WebDriver session
	caps := selenium.Capabilities{"browserName": "chrome"}
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", chromeDriverPort)) // ✅ FIXED URL
	if err != nil {
		log.Fatalf("Error starting WebDriver: %v", err)
	}
	defer driver.Quit()

	// Open The Verge Tech page
	url := "https://www.theverge.com/tech"
	if err := driver.Get(url); err != nil {
		log.Fatalf("Failed to load page: %v", err)
	}

	// Wait for page to load completely
	time.Sleep(5 * time.Second)

	// Extract article titles
	titles, err := driver.FindElements(selenium.ByCSSSelector, "a._1lkmsmo1")
	if err != nil {
		log.Fatalf("Error finding titles: %v", err)
	}

	// Extract article authors
	authors, err := driver.FindElements(selenium.ByCSSSelector, "span._1lldluw2 span")
	if err != nil {
		log.Fatalf("Error finding authors: %v", err)
	}

	// Extract article dates
	dates, err := driver.FindElements(selenium.ByCSSSelector, "span.duet--article--timestamp time")
	if err != nil {
		log.Fatalf("Error finding dates: %v", err)
	}

	// Extract article URLs
	urls, err := driver.FindElements(selenium.ByCSSSelector, "a._1lkmsmo1._184mftor")
	if err != nil {
		log.Fatalf("Error finding URLs: %v", err)
	}
	fmt.Println("Top 3 News Articles from The Verge:")

	for i := 0; i < 3; i++ {
		if i < len(titles) {
			title, _ := titles[i].Text()
			fmt.Printf("\n%d. Title: %s\n", i+1, title)
		}

		if i < len(authors) {
			author, _ := authors[i].Text()
			fmt.Printf("   Author: %s\n", author)
		}

		if i < len(dates) {
			date, _ := dates[i].Text()
			fmt.Printf("   Date: %s\n", date)
		}
		if i < len(urls) {
			url, _ := urls[i].GetAttribute("href")
			fmt.Printf("   URL: %s\n", url)
		}
	}
}