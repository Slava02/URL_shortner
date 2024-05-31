package app

import (
	"github.com/Slava02/URL_shortner/internal/app/handlers"
	"github.com/Slava02/URL_shortner/internal/app/storage"
	"net/http"
)

func Run() error {
	m := storage.NewUrlMapStorage()
	mux := http.NewServeMux()
	mux.Handle("/", &handlers.RootHandler{m})
	return http.ListenAndServe(`:8080`, mux)
}
