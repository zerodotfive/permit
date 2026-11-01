package activate

import (
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zerodotfive/permit/internal/app/types"
)

type Handler struct {
	postgres       *pgxpool.Pool
	logger         *log.Logger
	locationVerify string
	locationChain  string
}

func Register(pubRouter *http.ServeMux, container types.Container) {
	h := &Handler{
		postgres:       container.Postgres,
		logger:         container.Logger,
		locationVerify: container.Config.LocationVerify,
		locationChain:  container.Config.LocationChain,
	}

	url := "/" + strings.Trim(container.Config.LocationActivate, "/") + "/"
	pubRouter.HandleFunc("GET "+url+"{token}", h.get)
	pubRouter.HandleFunc("POST "+url+"{token}", h.post)
}
