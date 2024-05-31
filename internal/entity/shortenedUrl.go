package entity

import "github.com/Slava02/URL_shortner/internal/app/util"

type ShortenedURL struct {
	ID          string
	OriginalURL string
}

func NewShortenedURL(origUrl string) ShortenedURL {
	return ShortenedURL{ID: util.RandomString(), OriginalURL: origUrl}
}
