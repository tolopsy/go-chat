package main

import (
	"flag"
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

	fmt.Println("Request: ", r)
	t.templ.Execute(w, r)
}

func main() {
	addr := flag.String("addr", ":8080", "The port to listen from")
	flag.Parse()
	r := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	go r.run()

	log.Println("Starting web server on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
