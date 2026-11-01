package permits

import (
	"crypto/sha256"
)

func (p *Provider) VerifySession(sessionToken []byte) (*PermitSession, error) {
	sha := sha256.Sum256(sessionToken)

	var permit PermitSession
	row := p.postgres.QueryRow(p.ctx, `
		SELECT embed_site, target_site, session_token_use_not_after
		FROM permits
		WHERE session_token_sha = $1 AND session_token_use_not_after > now();`,
		sha[:],
	)

	err := row.Scan(&permit.EmbedSite, &permit.TargetSite, &permit.SessionTokenUseNotAfter)
	if err != nil {
		return nil, err
	}

	return &permit, nil
}
