package permits

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Provider struct {
	postgres *pgxpool.Pool
	ctx      context.Context
}

func NewProvider(ctx context.Context, postgres *pgxpool.Pool) *Provider {
	return &Provider{
		postgres: postgres,
		ctx:      ctx,
	}
}
