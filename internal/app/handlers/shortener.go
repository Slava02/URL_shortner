package handlers

import (
	"fmt"
	"github.com/Slava02/URL_shortner/internal/app/storage"
	"github.com/Slava02/URL_shortner/internal/app/util"
	"github.com/Slava02/URL_shortner/internal/config"
	"github.com/Slava02/URL_shortner/internal/entity"
	"io"
	"net/http"
	"strings"
)

type RootHandler struct {
	Storage storage.UrlStorage
}

func (r *RootHandler) post(w http.ResponseWriter, req *http.Request) {
	u, err := io.ReadAll(req.Body)
	origUrl := string(u)
	if err != nil || !util.IsURL(string(u)) {
		http.Error(w, fmt.Sprintf("Incorrect body:%Storage", origUrl), http.StatusBadRequest)
		return
	}

	shortUrl := entity.NewShortenedURL(origUrl)
	if err = r.Storage.Add(shortUrl); err != nil {
		http.Error(w, fmt.Sprintf("couldn't add shorturl %+v", shortUrl), http.StatusBadRequest)
		return
	}

	io.WriteString(w, fmt.Sprintf("%s/%s", config.Host, shortUrl.ID))
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
}

func (r *RootHandler) get(w http.ResponseWriter, req *http.Request) {
	rawPath := strings.Split(req.URL.Path, "/")
	id := rawPath[1]

	shortUrl, err := r.Storage.Get(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't get shorturl %w", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", shortUrl.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (r *RootHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.get(w, req)
	case http.MethodPost:
		r.post(w, req)
	default:
		http.Error(w, fmt.Sprintf("Incorrect method: %Storage", req.Method), http.StatusBadRequest)
	}
}
