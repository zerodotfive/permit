package unrestricted

import (
	"net/http"

	"github.com/zerodotfive/appcore/auth/handler"
	"github.com/zerodotfive/permit/internal/api/unrestricted/health"
	"github.com/zerodotfive/permit/internal/app/types"
)

func RegisterRoutes(apiRouter *http.ServeMux, container types.Container) {
	health.Register(apiRouter, container)
	handler.Register(apiRouter, container.Config.AuthConfig, container.Logger)
}
