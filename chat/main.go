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

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, req)
}

func main() {
	addr := flag.String("addr", ":8000", "The port to listen from")
	flag.Parse()

	r := createNewRoom()
	// r.tracer = trace.New(os.Stdout)

	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)
	fmt.Println("Running now")

	go r.run()

	log.Println("Running server at: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
