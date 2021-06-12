package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
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

	data := map[string]interface{}{
		"Host": req.Host,
	}

	if authCookie, err := req.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	addr := flag.String("addr", ":8000", "The port to listen from")
	flag.Parse()

	reader := JsonReader{}
	err := reader.reads()
	if err != nil {
		log.Fatal(err.Error())
	}

	gomniauth.SetSecurityKey(reader.ClientId)
	gomniauth.WithProviders(
		github.New(reader.ClientId,
			reader.ClientSecret,
			fmt.Sprintf("http://localhost%s/auth/callback/github", *addr)),
	)
	r := createNewRoom()
	// r.tracer = trace.New(os.Stdout)

	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	fmt.Println("Running now")

	go r.run()

	log.Println("Running server at: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
