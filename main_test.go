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

func TestUnauthorizedRequestReturns401(t *testing.T) {
	s := server{router: chi.NewRouter()}
	s.setupRoutes()

	r := httptest.NewRequest("GET", "/echo", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, r)
	if !reflect.DeepEqual(w.Code, http.StatusUnauthorized) {
		t.Fatalf("Response status not 401")
	}
}

func TestEchoApi(t *testing.T) {
	s := server{router: chi.NewRouter()}
	s.setupRoutes()

	expectedStr := "GET /echo HTTP/1.1\r"

	r := httptest.NewRequest("GET", "/echo", nil)
	r.Header.Set("Authorization", authKey)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, r)
	if !reflect.DeepEqual(w.Code, http.StatusOK) {
		t.Fatalf("Response status not 200")
	}

	responseStrs := strings.Split(w.Body.String(), "\n")

	if strings.Compare(expectedStr, responseStrs[0]) != 0 {
		fmt.Println("Expected: ", expectedStr)
		fmt.Println("Actual:   ", responseStrs[0])
		t.Fatalf("Expect response differs from actual")
	}
}

func TestEchoBodyApi(t *testing.T) {
	type body struct {
		Key string `json:"key"`
	}
	expectedBody := body{Key: "MyValue"}
	jsonBody, _ := json.Marshal(expectedBody)

	s := server{router: chi.NewRouter()}
	s.setupRoutes()

	r := httptest.NewRequest("GET", "/echoBody", bytes.NewReader(jsonBody))
	r.Header.Set("Authorization", authKey)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, r)
	if !reflect.DeepEqual(w.Code, http.StatusOK) {
		t.Fatalf("Response status not 200")
	}

	bodyBytes := w.Body.Bytes()
	responseBody := &body{}
	err := json.Unmarshal(bodyBytes, responseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal response")
	}

	if strings.Compare(expectedBody.Key, responseBody.Key) != 0 {
		fmt.Println("Expected: ", expectedBody.Key)
		fmt.Println("Actual:   ", responseBody.Key)
		t.Fatalf("Request and response body differ")
	}
}

type badBody struct{}

func (t badBody) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("error")
}

func TestEchoBodyReturns500OnCopyError(t *testing.T) {
	s := server{router: chi.NewRouter()}
	s.setupRoutes()

	r := httptest.NewRequest("GET", "/echoBody", badBody{})
	r.Header.Set("Authorization", authKey)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, r)
	if !reflect.DeepEqual(w.Code, http.StatusInternalServerError) {
		t.Fatalf("Response status not 500")
	}
}
