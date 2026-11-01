package admin

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/zerodotfive/permit/internal/app/types"
)

func RegisterRoutes(hapi huma.API, container types.Container) {
	Register(hapi, container)
}
