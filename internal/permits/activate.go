package permits

import (
	"crypto/rand"
	"crypto/sha256"
	"uuid"

	"github.com/jackc/pgx/v5"
)

func (p *Provider) ActivateToken(token []byte) (*PermitSession, error) {
	tokenSha := sha256.Sum256(token)

	sessionToken := make([]byte, 32)
	if _, err := rand.Read(sessionToken); err != nil {
		return nil, err
	}

	sessionTokenSha := sha256.Sum256(sessionToken)

	chainToken := make([]byte, 32)
	if _, err := rand.Read(chainToken); err != nil {
		return nil, err
	}

	chainTokenSha := sha256.Sum256(chainToken)

	row := p.postgres.QueryRow(
		p.ctx, `
			WITH parent AS (
				UPDATE permits
					SET session_token_sha = @session_token_sha, session_token_use_not_after = now() + make_interval(mins => session_lifetime_minutes)
				WHERE
					token_sha = @token_sha AND parent_id IS NULL AND session_token_sha IS NULL AND token_use_not_after > now()
				RETURNING
					id, comment, embed_site, target_site, target_location, session_token_use_not_after
				), 
			child AS (
				INSERT INTO permits
					(id, parent_id, comment, embed_site, target_site, target_location, token_sha, token_use_not_after, session_lifetime_minutes)
				SELECT
					@chain_id, p.id, p.comment, p.embed_site, p.target_site, p.target_location, @chain_token_sha, now() + make_interval(mins => 1), 0
				FROM parent p
			)
			SELECT embed_site, target_site, target_location, session_token_use_not_after FROM parent p;
		`,
		pgx.NamedArgs{
			"session_token_sha": sessionTokenSha[:],
			"token_sha":         tokenSha[:],
			"chain_id":          uuid.New(),
			"chain_token_sha":   chainTokenSha[:],
		},
	)

	permit := &PermitSession{
		SessionToken: sessionToken,
		ChainToken:   chainToken,
	}
	err := row.Scan(
		&permit.EmbedSite,
		&permit.TargetSite,
		&permit.TargetLocation,
		&permit.SessionTokenUseNotAfter,
	)
	if err != nil {
		return nil, err
	}

	return permit, nil
}
