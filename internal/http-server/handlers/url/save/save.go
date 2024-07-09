package save

import (
	"errors"
	resp "github.com/Slava02/URL_shortner/internal/lib/api/response"
	"github.com/Slava02/URL_shortner/internal/lib/random"
	"github.com/Slava02/URL_shortner/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config
const aliasLength = 4
const maxRetries = 5

//go:generate go run .github.com/vektra/mockery/v2@v2.43.2 --name=URLSaver
type URLSaver interface {
	SaveUrl(urlToSave string, alias string) (int64, error)
	ExistsAlias(alias string) (bool, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))

			render.JSON(w, r, resp.Error("failed to decode request body"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.Any("error", err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		log.Info("request body validated", slog.Any("request", req))

		alias := req.Alias
		if alias == "" {
			retries := 1
			for retries <= maxRetries {
				alias = random.NewRandomString(aliasLength)
				if exists, err := urlSaver.ExistsAlias(alias); !exists && err == nil {
					break
				} else if err != nil {
					log.Error("can't check alias existence", slog.Any("retrie: ", retries), slog.Any("error", err))
				}
				retries++
			}
			if retries == maxRetries {
				log.Error("couldn't check alias existence", slog.Any("error", err))
			} else {
				log.Info("alias generated", slog.Any("request", req))
			}
		}

		id, err := urlSaver.SaveUrl(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", slog.Any("error", err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
