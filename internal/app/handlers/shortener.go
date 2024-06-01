package handlers

import (
	"fmt"
	"github.com/Slava02/URL_shortner/internal/app/storage"
	"github.com/Slava02/URL_shortner/internal/app/util"
	"github.com/Slava02/URL_shortner/internal/config"
	"github.com/Slava02/URL_shortner/internal/entity"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type RootHandler struct {
	Storage storage.UrlStorage
}

func (r *RootHandler) post(w http.ResponseWriter, req *http.Request) {
	u, err := io.ReadAll(req.Body)
	origUrl := string(u)
	if err != nil || !util.IsURL(string(u)) {
		http.Error(w, fmt.Sprintf("Incorrect body:%s", origUrl), http.StatusBadRequest)
		return
	}

	shortUrl := entity.NewShortenedURL(origUrl)
	if err = r.Storage.Add(shortUrl); err != nil {
		http.Error(w, fmt.Sprintf("couldn't add shorturl %+v", shortUrl), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, fmt.Sprintf("%s/%s", config.Host, shortUrl.ID))
}

func (r *RootHandler) get(w http.ResponseWriter, req *http.Request) {
	rawPath := strings.Split(req.URL.Path, "/")
	id := rawPath[1]

	shortUrl, err := r.Storage.Get(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't get shorturl %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", shortUrl.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (r *RootHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pat := regexp.MustCompile(`[A-Z,a-z,0-9]{8}`)
	switch q := strings.TrimPrefix(req.URL.Path, "/"); {
	case q == "" && req.Method == http.MethodPost:
		r.post(w, req)
	case pat.MatchString(q) && req.Method == http.MethodGet:
		r.get(w, req)
	default:
		http.Error(w, fmt.Sprintf("Incorrect method: %s", req.Method), http.StatusBadRequest)
	}

}
