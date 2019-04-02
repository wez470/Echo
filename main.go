package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	"io"
	"log"
	"net/http"
	"strings"
)

const port = "8080"
const authKey = "Basic dXNlcjpwYXNzd29yZA=="

var generateDocs = flag.Bool("docs", false, "Generate server route documentation")

//server represents the application server
type server struct {
	router *chi.Mux
}

//setupRoutes sets up application routes and middleware
func (s *server) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.RedirectSlashes)
	s.router.Use(authenticationMiddleware)
	s.router.HandleFunc("/echoBody", s.echoBody())
	s.router.HandleFunc("/echo", s.echo())
}

func authenticationMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if len(authHeader) != 1 || strings.Compare(authHeader[0], authKey) != 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized."))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//echoBody returns a handler that returns the request body
func (s *server) echoBody() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(w, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to read request body."))
			return
		}
	}
}

//echo returns a handler that returns the request
func (s *server) echo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.Write(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to write response."))
		}
	}
}

//printAPIDocs prints out API docs as JSON
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
