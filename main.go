package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// handles a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// Ensures that template is only parsed once and then executes template
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, nil)
}

func main() {
	r := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/", r)

	go r.run()

	fmt.Println("Running")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
