package main

import (
	"github.com/spf13/viper"
	appcoreAuthTypes "github.com/zerodotfive/appcore/auth/types"
	"github.com/zerodotfive/permit/internal/app"
	"github.com/zerodotfive/permit/internal/app/types"
)

func main() {
	viper.SetDefault("PUB_LISTEN", "127.0.0.1:8090")
	viper.SetDefault("PUB_LOCATION_ACTIVATE", "/.p/a")
	viper.SetDefault("PUB_LOCATION_CHAIN", "/.p/c")
	viper.SetDefault("PUB_LOCATION_VERIFY", "/.p/v")
	viper.SetDefault("API_LISTEN", "127.0.0.1:8091")
	viper.SetDefault("DOCS_ENABLED", false)
	viper.SetDefault("AUTH_ENABLED", true)
	viper.SetDefault("AUTH_SCOPES", "openid profile email groups")
	viper.SetDefault("POSTGRES_URI", "postgres://postgres:password@localhost:5432/permit?sslmode=prefer")
	viper.SetDefault("MIGRATE", true)
	viper.AutomaticEnv()

	app.Run(&types.Config{
		PubListen:        viper.GetString("PUB_LISTEN"),
		LocationActivate: viper.GetString("PUB_LOCATION_ACTIVATE"),
		LocationChain:    viper.GetString("PUB_LOCATION_CHAIN"),
		LocationVerify:   viper.GetString("PUB_LOCATION_VERIFY"),
		ApiListen:        viper.GetString("API_LISTEN"),
		DocsEnabled:      viper.GetBool("DOCS_ENABLED"),
		AuthEnabled:      viper.GetBool("AUTH_ENABLED"),
		AuthConfig: &appcoreAuthTypes.AuthConfig{
			AuthClientID:              viper.GetString("AUTH_CLIENT_ID"),
			AuthClientSecret:          viper.GetString("AUTH_CLIENT_SECRET"),
			AuthRedirectURL:           viper.GetString("AUTH_REDIRECT_URL"),
			AuthScopes:                viper.GetStringSlice("AUTH_SCOPES"),
			OpenIDConfigurationIssuer: viper.GetString("OPENID_CONFIGURATION_ISSUER"),
			JWTKey:                    viper.GetString("JWT_KEY"),
			AdminGroup:                viper.GetString("ADMIN_GROUP"),
		},
		PostgresURI: viper.GetString("POSTGRES_URI"),
		Migrate:     viper.GetBool("MIGRATE"),
	})
}
