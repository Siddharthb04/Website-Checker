package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

// PageData holds the data for rendering
type PageData struct {
	Results map[string]string
}

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func checkStatus(urls []string) map[string]string {
	results := make(map[string]string)
	client := http.Client{Timeout: 5 * time.Second}

	for _, url := range urls {
		resp, err := client.Get(url)
		if err != nil {
			results[url] = "❌ Offline"
		} else {
			results[url] = "✅ " + http.StatusText(resp.StatusCode)
			resp.Body.Close()
		}
	}
	return results
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		urls := r.Form["urls"]
		data := PageData{Results: checkStatus(urls)}
		tmpl.Execute(w, data)
		return
	}

	tmpl.Execute(w, PageData{})
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
