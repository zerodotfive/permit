package health

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zerodotfive/permit/internal/app/types"
)

type Handler struct {
	postgres *pgxpool.Pool
	logger   *log.Logger
}

func Register(apiRouter *http.ServeMux, container types.Container) {
	h := &Handler{
		postgres: container.Postgres,
		logger:   container.Logger,
	}

	apiRouter.HandleFunc("GET /health", h.get)
}
