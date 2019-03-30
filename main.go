package main

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//Server the application server
type server struct{}

func (s *server) echoBody() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf, _ := ioutil.ReadAll(r.Body)
		w.Write(buf)
	}
}

func (s *server) echo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Write(w)
	}
}

func main() {
	s := server{}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RedirectSlashes)
	r.HandleFunc("/echo", s.echoBody())
	r.HandleFunc("/echoAll", s.echo())
	http.ListenAndServe(":8080", r)
}
