package permits

import (
	"crypto/rand"
	"crypto/sha256"

	"github.com/jackc/pgx/v5"
)

func (p *Provider) ChainToken(token []byte) (*PermitSession, error) {
	tokenSha := sha256.Sum256(token)
	sessionToken := make([]byte, 32)
	if _, err := rand.Read(sessionToken); err != nil {
		return nil, err
	}

	sessionTokenSha := sha256.Sum256(sessionToken)

	row := p.postgres.QueryRow(
		p.ctx,
		`
			UPDATE permits AS c
			SET session_token_sha = @session_token_sha,
				session_token_use_not_after = p.session_token_use_not_after
			FROM permits AS p
			WHERE c.token_sha = @token_sha
			  AND c.parent_id IS NOT NULL
			  AND p.id = c.parent_id
			  AND p.session_token_use_not_after > now()
			  AND c.session_token_sha IS NULL
			  AND c.token_use_not_after > now()
			RETURNING c.embed_site, c.target_site, c.target_location, c.session_token_use_not_after;`,
		pgx.NamedArgs{
			"session_token_sha": sessionTokenSha[:],
			"token_sha":         tokenSha[:],
		},
	)

	permit := &PermitSession{
		SessionToken: sessionToken,
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
