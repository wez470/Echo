package main

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const port = "8080"

//Server the application server
type server struct{}

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

func main() {
	s := server{}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RedirectSlashes)
	r.HandleFunc("/echoBody", s.echoBody())
	r.HandleFunc("/echo", s.echo())
	http.ListenAndServe(":"+port, r)
}
