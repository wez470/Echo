package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestEchoApi(t *testing.T) {
	s := server{router: chi.NewRouter()}
	s.setupRoutes()

	expectedStrs := []string{"GET /echo HTTP/1.1\r", "Host: \r", "User-Agent: Go-http-client/1.1\r", "\r", ""}

	r, err := http.NewRequest("GET", "/echo", nil)
	if err != nil {
		t.Fatalf("Failed to create fake http request")
	}
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, r)
	if !reflect.DeepEqual(w.Code, http.StatusOK) {
		t.Fatalf("Response status not 200")
	}

	responseStrs := strings.Split(w.Body.String(), "\n")

	for i := range responseStrs {
		if strings.Compare(expectedStrs[i], responseStrs[i]) != 0 {
			fmt.Println("Expected: ", expectedStrs[i])
			fmt.Println("Actual:   ", responseStrs[i])
			t.Fatalf("Expect response differs from actual")
		}
	}
}

func TestEchoBodyApi(t *testing.T) {
	type body struct {
		Key string `json:"key"`
	}
	expectedBody := body{Key: "MyValue"}
	jsonBody, err := json.Marshal(expectedBody)

	s := server{router: chi.NewRouter()}
	s.setupRoutes()

	r, err := http.NewRequest("GET", "/echoBody", bytes.NewReader(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create fake http request")
	}
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, r)
	if !reflect.DeepEqual(w.Code, http.StatusOK) {
		t.Fatalf("Response status not 200")
	}

	bodyBytes := w.Body.Bytes()
	responseBody := &body{}
	err = json.Unmarshal(bodyBytes, responseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal response")
	}

	if strings.Compare(expectedBody.Key, responseBody.Key) != 0 {
		fmt.Println("Expected: ", expectedBody.Key)
		fmt.Println("Actual:   ", responseBody.Key)
		t.Fatalf("Request and response body differ")
	}
}