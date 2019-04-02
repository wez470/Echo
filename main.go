package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
)

const port = "8080"

var generateDocs = flag.Bool("docs", false, "Generate server route documentation")

//server represents the application server
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
		_, err := io.Copy(w, r.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Failed to read request body."))
			return
		}
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

func (s *server) printAPIDocs() {
	fmt.Println(docgen.JSONRoutesDoc(s.router))
}

func main() {
	flag.Parse()
	s := server{router: chi.NewRouter()}
	s.setupRoutes()
	if *generateDocs {
		s.printAPIDocs()
		return
	}
	log.Printf("Starting server on port: %v\n", port)
	http.ListenAndServe(":"+port, s.router)
}
