package del

import (
	resp "github.com/Slava02/URL_shortner/internal/lib/api/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLDeleter
type URLDeleter interface {
	DeleteUrl(id int) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.del.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id_param := chi.URLParam(r, "id")
		if id_param == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		id, err := strconv.Atoi(id_param)
		if err != nil {
			log.Error("failed to convert id", slog.Any("error", err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		err = urlDeleter.DeleteUrl(id)
		if err != nil {
			log.Error("failed to delete url", slog.Any("error", err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("deleted url", slog.String("id", id_param))

		render.JSON(w, r, resp.OK())
	}
}
