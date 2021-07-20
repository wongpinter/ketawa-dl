package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"time"
)

type Post struct {
	Title			string
	Category		string
	PublishedAt		string
	Content			string
}

func main() {
	posts := make([]Post, 0, 20)

	c := colly.NewCollector(
		colly.Async(true),
	)

	contentCollector := c.Clone()

	c.OnHTML(`div.row.menu > ul > li > a`, func(e *colly.HTMLElement) {
		categoryUrl := e.Request.AbsoluteURL(e.Attr("href"))

		_ = c.Visit(categoryUrl)
	})

	c.OnHTML(`p.cont-page:last-child > a:last-child`, func(e *colly.HTMLElement) {
		pageUrl := e.Request.AbsoluteURL(e.Attr("href"))

		_ = contentCollector.Visit(pageUrl)
	})


	contentCollector.OnHTML(`ul.listdata > li:last-child > a`, func(e *colly.HTMLElement) {
		_ = contentCollector.Visit(e.Attr(`href`))
	})

	contentCollector.OnHTML("div#prevnext > a", func(e *colly.HTMLElement) {
		pageLink := e.Request.AbsoluteURL(e.Attr("href"))

		_ = contentCollector.Visit(pageLink)
	})

	contentCollector.OnHTML("div.post", func(e *colly.HTMLElement) {
		temp := Post{}
		temp.Category = e.ChildText("div.categories > span:nth-child(2) > a")
		temp.Title = e.ChildText("h2.title")
		temp.Content = e.ChildText(".entry")
		temp.PublishedAt = e.ChildText("p.meta > span:nth-child(2)")

		posts = append(posts, temp)
	})

	contentCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	_ = c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	_ = c.Visit("https://www.ketawa.com/")

	c.Wait()
	contentCollector.Wait()

	file, _ := json.MarshalIndent(posts, "", " ")
	_ = ioutil.WriteFile("ketawa.json", file, 0644)

	fmt.Println("Scraping Done!..")
}
