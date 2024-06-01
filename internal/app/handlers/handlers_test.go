package handlers

import (
	"bytes"
	"github.com/Slava02/URL_shortner/internal/app/storage"
	"github.com/Slava02/URL_shortner/internal/entity"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestRootHandler_post(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		method  string
		request string
		body    string
		want    want
	}{
		{
			name:   "simple post test #1",
			method: http.MethodPost,
			body:   "https://practicum.yandex.ru/",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
			},
			request: "/",
		},
		{
			name:   "no body post test #2",
			method: http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
			request: "/",
		},
		{
			name:   "prohibited request for post test #3",
			method: http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
			request: "/fdaw43",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewUrlMapStorage()
			root := RootHandler{store}

			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.request, body)
			w := httptest.NewRecorder()

			root.post(w, request)

			result := w.Result()
			buf := new(bytes.Buffer)
			buf.ReadFrom(result.Body)
			shortStr := buf.String()

			patterStr := tt.want.contentType + "*"
			pat := regexp.MustCompile(patterStr)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.True(t, pat.MatchString(result.Header.Get("Content-Type")))
			if result.StatusCode != http.StatusBadRequest {
				assert.Len(t, shortStr, 30)
			}
		})
	}
}

var testShortUrl = map[string]entity.ShortenedURL{
	"simple": {
		ID:          "12345678",
		OriginalURL: "ya.ru",
	},
}

func TestRootHandler_get(t *testing.T) {
	type want struct {
		statusCode int
		url        string
	}
	tests := []struct {
		name    string
		method  string
		request string
		want    want
	}{
		{
			name:   "simple get test #1",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				url:        testShortUrl["simple"].OriginalURL,
			},
			request: "/" + testShortUrl["simple"].ID,
		},
		{
			name:   "not allowed id get test #2",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
			},
			request: "/not-allowed-id",
		},
		{
			name:   "root request get test #2",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
			},
			request: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewUrlMapStorage()
			root := RootHandler{store}
			store.Add(testShortUrl["simple"])

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			root.get(w, request)

			result := w.Result()

			if tt.want.url != "" {
				assert.Equal(t, tt.want.url, result.Header.Get("Location"))
			}
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
