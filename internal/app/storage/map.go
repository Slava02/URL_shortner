package storage

import (
	"fmt"
	"github.com/Slava02/URL_shortner/internal/entity"
)

type UrlMapStorage struct {
	shortenedUrl map[string]string
}

func NewUrlMapStorage() UrlMapStorage {
	return UrlMapStorage{make(map[string]string)}
}

func (m UrlMapStorage) Add(u entity.ShortenedURL) error {
	_, ok := m.shortenedUrl[u.ID]
	if ok {
		return fmt.Errorf("can't add: key already exists")
	}
	m.shortenedUrl[u.ID] = u.OriginalURL
	return nil
}

func (m UrlMapStorage) Get(id string) (*entity.ShortenedURL, error) {
	u, ok := m.shortenedUrl[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &entity.ShortenedURL{ID: id, OriginalURL: u}, nil
}
