package types

import "github.com/zerodotfive/appcore/auth/types"

type Config struct {
	PubListen        string
	LocationActivate string
	LocationChain    string
	LocationVerify   string
	ApiListen        string
	DocsEnabled      bool
	AuthEnabled      bool
	PostgresURI      string
	Migrate          bool
	AuthConfig       *types.AuthConfig
}
