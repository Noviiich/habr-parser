package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type Articles struct {
	name, title, description string
}

func main() {
	var articles []Articles

	service, err := selenium.NewChromeDriverService("./chromedriver/chromedriver.exe", 4444)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--headless-new",
	}})

	driver, err := selenium.NewRemote(caps, "")

	if err != nil {
		log.Fatal("Error:", err)
	}

	err = driver.MaximizeWindow("")
	if err != nil {
		log.Fatal("Error:", err)
	}

	err = driver.Get("https://habr.com/ru/articles/")
	if err != nil {
		log.Fatal("Error:", err)

		html, err := driver.PageSource()
		if err != nil {
			log.Fatal("Error:", err)
		}
		fmt.Println(html)
	}

	nameArticles, err := driver.FindElements(selenium.ByCSSSelector, ".tm-user-info__user")
	if err != nil {
		log.Fatal("Error:", err)
	}

	titleArticles, err := driver.FindElements(selenium.ByCSSSelector, ".tm-title__link")
	if err != nil {
		log.Fatal("Error:", err)
	}

	descriptionArticles, err := driver.FindElements(selenium.ByCSSSelector, ".article-formatted-body_version-2")
	if err != nil {
		log.Fatal("Error:", err)
	}

	for i := 0; i < 16; i++ {
		nameArticle := nameArticles[i]
		titleArticle := titleArticles[i]
		descriptionArticle := descriptionArticles[i]
		nameElement, err := nameArticle.FindElement(selenium.ByCSSSelector, "a")
		titleElement, err := titleArticle.FindElement(selenium.ByCSSSelector, "span")
		descriptionElement, err := descriptionArticle.FindElement(selenium.ByCSSSelector, "p")
		name, err := nameElement.Text()
		title, err := titleElement.Text()
		description, err := descriptionElement.Text()
		if err != nil {
			log.Fatal("Error:", err)
		}

		article := Articles{}
		article.name = name
		article.title = title
		article.description = description
		articles = append(articles, article)
	}

	fmt.Println(articles)

	file, err := os.Create("articles.csv")
	if err != nil {
		log.Fatal("Error:", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"name",
		"title",
		"description",
	}

	writer.Write(headers)

	for _, article := range articles {
		record := []string{
			article.name,
			article.title,
			article.description,
		}

		writer.Write(record)
	}
	defer writer.Flush()
}
