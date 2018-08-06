package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
)

type pageInfo struct {
	StatusCode int
	Links      map[string]string
}

var pages map[string]pageInfo

func parseH(url string) pageInfo {
	log.Println("visiting", url)

	c := colly.NewCollector()

	p := &pageInfo{Links: make(map[string]string)}

	// count links
	count := 1
	c.OnHTML(".itemlist tr.athing td.title", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		title := e.ChildText("a")
		span := e.ChildText(".sitestr")
		if link != "" && count <= 10 {
			p.Links[link] = strings.Split(title, span)[0]
			count++
		}
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		p.StatusCode = r.StatusCode
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	c.Visit(url)

	return *p
}

func parseHN(url string) pageInfo {
	log.Println("visiting", url)

	c := colly.NewCollector()

	p := &pageInfo{Links: make(map[string]string)}

	// count links
	count := 1
	c.OnHTML(".js-trackedPost", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		title := e.ChildText("h3")
		if link != "" && count <= 10 {
			p.Links[link] = title
			count++
		}
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		p.StatusCode = r.StatusCode
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	c.Visit(url)

	return *p
}

func parseTC(url string) pageInfo {
	log.Println("visiting", url)

	c := colly.NewCollector()

	p := &pageInfo{Links: make(map[string]string)}

	// count links
	count := 1
	c.OnHTML(".post-block__title", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		title := e.ChildText("a")
		if link != "" && count <= 10 {
			p.Links[link] = title
			count++
		}
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		p.StatusCode = r.StatusCode
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	c.Visit(url)

	return *p
}

func handler(w http.ResponseWriter, r *http.Request) {

	pages := make(map[string]pageInfo)

	url1 := "https://news.ycombinator.com/"
	url2 := "https://hackernoon.com/"
	url3 := "https://techcrunch.com/"

	p1 := parseH(url1)
	p2 := parseHN(url2)
	p3 := parseTC(url3)

	pages["Hackernews"] = p1
	pages["HackerNoon"] = p2
	pages["TechCrunch"] = p3

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, pages)
}

func main() {
	addr := ":8080"

	http.HandleFunc("/", handler)

	log.Println("listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
