package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

type hackernewsArticles struct {
	title string
	link  string
}

type pageInfo struct {
	StatusCode int
	Links      map[string]int
}

func handler(w http.ResponseWriter, r *http.Request) {
	URL := "https://news.ycombinator.com/"
	if URL == "" {
		log.Println("missing URL argument")
		return
	}
	log.Println("visiting", URL)

	c := colly.NewCollector()

	p := &pageInfo{Links: make(map[string]int)}

	// count links
	c.OnHTML(".itemlist tr.athing td.title", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		log.Println(e.ChildText("a"))
		if link != "" {
			p.Links[link]++
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

	c.Visit(URL)

	// dump results
	b, err := json.Marshal(p)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, "")
}

func main() {
	addr := ":8080"

	http.HandleFunc("/", handler)
	http.HandleFunc("/test", templateHandler)

	log.Println("listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
