package storage

import "github.com/Slava02/URL_shortner/internal/entity"

type UrlStorage interface {
	Add(entity.ShortenedURL) error
	Get(id string) (*entity.ShortenedURL, error)
}
