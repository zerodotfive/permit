package pub

import (
	"net/http"

	"github.com/zerodotfive/permit/internal/app/types"
	"github.com/zerodotfive/permit/internal/pub/activate"
	"github.com/zerodotfive/permit/internal/pub/chain"
	"github.com/zerodotfive/permit/internal/pub/verify"
)

func RegisterRoutes(pubRouter *http.ServeMux, container types.Container) {
	activate.Register(pubRouter, container)
	chain.Register(pubRouter, container)
	verify.Register(pubRouter, container)
}
