package permits

import (
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zerodotfive/permit/internal/app/types"
)

type Handler struct {
	postgres         *pgxpool.Pool
	logger           *log.Logger
	locationActivate string
}

func Register(hapi huma.API, container types.Container) {
	h := &Handler{
		postgres:         container.Postgres,
		logger:           container.Logger,
		locationActivate: container.Config.LocationActivate,
	}

	huma.Register(hapi, huma.Operation{
		OperationID: "permit-create",
		Method:      http.MethodPost,
		Path:        "/permits",
		Summary:     "Create permit",
		Tags:        []string{"Permit"},
	}, h.create)

	huma.Register(hapi, huma.Operation{
		OperationID: "permit-list",
		Method:      http.MethodGet,
		Path:        "/permits",
		Summary:     "List permits",
		Tags:        []string{"Permit"},
	}, h.list)

	huma.Register(hapi, huma.Operation{
		OperationID: "permit-delete",
		Method:      http.MethodDelete,
		Path:        "/permits/{id}",
		Summary:     "Delete permit",
		Tags:        []string{"Permit"},
	}, h.delete)
}
