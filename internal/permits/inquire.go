package permits

import (
	"crypto/sha256"
)

func (p *Provider) InquireForToken(token []byte) (*PermitSession, error) {
	tokenSha := sha256.Sum256(token)

	row := p.postgres.QueryRow(
		p.ctx,
		`
			SELECT embed_site, target_site, session_token_use_not_after
			FROM permits
			WHERE token_sha = $1 AND session_token_sha IS NULL AND token_use_not_after > now();`,
		tokenSha[:],
	)

	permit := &PermitSession{}
	err := row.Scan(
		&permit.EmbedSite,
		&permit.TargetSite,
		&permit.SessionTokenUseNotAfter,
	)
	if err != nil {
		return nil, err
	}

	return permit, nil
}
