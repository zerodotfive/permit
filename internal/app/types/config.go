package types

import "github.com/zerodotfive/appcore/auth/oidc"

type Config struct {
	PubListen          string
	LocationActivate   string
	LocationChain      string
	LocationVerify     string
	ApiListen          string
	DocsEnabled        bool
	AuthEnabled        bool
	PostgresURI        string
	Migrate            bool
	OAuthConfig        *oidc.OAuthConfig
	AdminGroup         string
	StaticTokenEnabled bool
	StaticTokenSHA256  string
}
