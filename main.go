package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
)

const port = "8080"

var generateDocs = flag.Bool("docs", false, "Generate server route documentation")

//Server represents the application server
type server struct {
	router *chi.Mux
}

func (s *server) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.RedirectSlashes)
	s.router.HandleFunc("/echoBody", s.echoBody())
	s.router.HandleFunc("/echo", s.echo())
}

func (s *server) echoBody() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Failed to read request body."))
			return
		}
		w.Write(buf)
	}
}

func (s *server) echo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.Write(w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Failed to write response."))
		}
	}
}

func (s *server) conditionallyPrintAPIDocs() {
	if *generateDocs {
		fmt.Println(docgen.MarkdownRoutesDoc(s.router, docgen.MarkdownOpts{
			ProjectPath: "github.com/wez470/Echo",
			Intro:       "Echo server generated docs.",
		}))
		return
	}
}

func main() {
	flag.Parse()
	s := server{router: chi.NewRouter()}
	s.setupRoutes()
	s.conditionallyPrintAPIDocs()
	log.Printf("Starting server on port: %v\n", port)
	http.ListenAndServe(":"+port, s.router)
}
