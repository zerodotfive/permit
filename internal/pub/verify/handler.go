package verify

import (
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zerodotfive/permit/internal/app/types"
)

type Handler struct {
	postgres *pgxpool.Pool
	logger   *log.Logger
}

func Register(pubRouter *http.ServeMux, container types.Container) {
	h := &Handler{
		postgres: container.Postgres,
		logger:   container.Logger,
	}

	url := "/" + strings.Trim(container.Config.LocationVerify, "/")
	pubRouter.HandleFunc("GET "+url, h.get)
}
