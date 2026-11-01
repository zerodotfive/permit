package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	appcoreAuthHelper "github.com/zerodotfive/appcore/auth/helper"
	appcoreAuthTypes "github.com/zerodotfive/appcore/auth/types"
	"github.com/zerodotfive/permit/internal/api/admin"
	"github.com/zerodotfive/permit/internal/api/unrestricted"
	"github.com/zerodotfive/permit/internal/app/types"
	"github.com/zerodotfive/permit/internal/info"
	"github.com/zerodotfive/permit/internal/permits"
	"github.com/zerodotfive/permit/internal/pub"
	adminUI "github.com/zerodotfive/permit/internal/ui/admin"
	"github.com/zerodotfive/permit/migrations"
	"github.com/zerodotfive/permit/ui"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func Run(config *types.Config) {
	log.SetPrefix("[" + info.Application + "] ")
	log.Printf(info.Application+" version %s", info.ApplicationVersion)

	if config.AuthEnabled && len(config.AuthConfig.JWTKey) < 32 {
		log.Fatal(errors.New("JWT key length must be at least 32 bytes"))
	}

	postgres, err := pgxpool.New(context.Background(), config.PostgresURI)
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.Close()

	container := types.Container{
		Config:   config,
		Logger:   log.Default(),
		Postgres: postgres,
	}

	if config.Migrate {
		db := stdlib.OpenDBFromPool(postgres)
		defer func() { _ = db.Close() }()

		goose.SetBaseFS(migrations.EmbedMigrations)

		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatal(err)
		}

		if err := goose.Up(db, "."); err != nil {
			log.Fatal(err)
		}
	}

	go permits.RunCleanup(postgres, container.Logger, 30*time.Minute, 24*time.Hour)

	humaConfig := huma.DefaultConfig(info.Application, info.ApplicationVersion)
	humaConfig.DocsRenderer = huma.DocsRendererStoplightElements

	if !config.DocsEnabled {
		humaConfig.DocsRenderer = ""
		humaConfig.DocsPath = ""
		humaConfig.OpenAPIPath = ""
		humaConfig.SchemasPath = ""
	}

	unauthorized, _ := ui.GetEmbedHTMLBytes("unauthorized.html")
	forbidden, _ := ui.GetEmbedHTMLBytes("forbidden.html")
	statusPages := appcoreAuthTypes.StatusPages{
		Unauthorized: unauthorized,
		Forbidden:    forbidden,
	}

	pubRouter := http.NewServeMux()
	apiRouter := http.NewServeMux()

	pub.RegisterRoutes(pubRouter, container)

	unrestricted.RegisterRoutes(apiRouter, container)

	hapi := humago.New(apiRouter, humaConfig)
	humaRoot := huma.NewGroup(hapi, "")

	humaAdminRestricted := huma.NewGroup(humaRoot, "")
	if config.AuthEnabled {
		humaAdminRestricted.UseMiddleware(appcoreAuthHelper.NewAdminAuthMiddleware(container.Config.AuthConfig, &statusPages))
	}

	humaAdminRestrictedAPI := huma.NewGroup(humaAdminRestricted, "/api/"+info.ApiVersion)
	admin.RegisterRoutes(humaAdminRestrictedAPI, container)

	humaAdminRestrictedUI := huma.NewGroup(humaAdminRestricted, "/")
	adminUI.RegisterRoutes(humaAdminRestrictedUI, container)

	errCh := make(chan error, 2)

	go func() {
		errCh <- http.ListenAndServe(config.PubListen, pubRouter)
	}()

	go func() {
		errCh <- http.ListenAndServe(config.ApiListen, apiRouter)
	}()

	if err := <-errCh; err != nil {
		log.Fatal(err)
	}
}
