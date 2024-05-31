package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

var urls = make(map[string]string)
var urlLen = 8
var host = "http://localhost:8080/"

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", MainPage)
	return http.ListenAndServe(`:8080`, mux)
}

func MainPage(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = url.ParseRequestURI(string(body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		shortenedUrl := UrlShortner()
		urls[shortenedUrl] = string(body)
		res := host + shortenedUrl
		fmt.Fprintf(w, res)
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
	} else if req.Method == http.MethodGet {
		log.Printf("Path:%s\nRawPath:%s\n\nMap:%v\n", req.URL.Path, string(bytes.TrimPrefix([]byte(req.URL.Path), []byte("/"))), urls)
		u, ok := urls[string(bytes.TrimPrefix([]byte(req.URL.Path), []byte("/")))]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("Location", u)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	return
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func UrlShortner() string {
	b := make([]byte, urlLen)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
