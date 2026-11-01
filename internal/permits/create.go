package permits

import (
	"crypto/rand"
	"crypto/sha256"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
)

func (p *Provider) Create(comment string, embedSite string, targetSite string, targetLocation string, tokenLifetimeMinutes int, sessionLifetimeMinutes int) (*PermitToken, error) {
	if tokenLifetimeMinutes < 1 || tokenLifetimeMinutes > 1440 {
		return nil, InvalidTokenLifetime
	}

	if sessionLifetimeMinutes < 1 || sessionLifetimeMinutes > 1440 {
		return nil, InvalidSessionLifetime
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return nil, err
	}

	tokenUseNotAfter := time.Now().Add(time.Minute * time.Duration(tokenLifetimeMinutes))

	tokenSHA := sha256.Sum256(token)

	permit := PermitToken{
		ID:               uuid.New(),
		Comment:          comment,
		EmbedSite:        embedSite,
		TargetSite:       targetSite,
		TargetLocation:   targetLocation,
		Token:            token,
		TokenUseNotAfter: tokenUseNotAfter,
	}

	_, err := p.postgres.Exec(p.ctx, `
		INSERT INTO permits
		(id, comment, embed_site, target_site, target_location, token_sha, token_use_not_after, session_lifetime_minutes)
		VALUES
		(@id, @comment, @embed_site, @target_site, @target_location, @token_sha, @token_use_not_after, @session_lifetime_minutes)`,
		pgx.NamedArgs{
			"id":                       permit.ID,
			"comment":                  permit.Comment,
			"embed_site":               permit.EmbedSite,
			"target_site":              permit.TargetSite,
			"target_location":          permit.TargetLocation,
			"token_sha":                tokenSHA[:],
			"token_use_not_after":      permit.TokenUseNotAfter,
			"session_lifetime_minutes": sessionLifetimeMinutes,
		},
	)
	if err != nil {
		return nil, err
	}

	return &permit, nil
}
