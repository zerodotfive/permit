package permits

func (p *Provider) List() ([]PermitShort, error) {
	rows, err := p.postgres.Query(
		p.ctx,
		`SELECT id, parent_id, comment, embed_site, target_site, target_location, token_use_not_after, session_token_use_not_after
		FROM permits;`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permits []PermitShort
	for rows.Next() {
		var permit PermitShort
		err := rows.Scan(
			&permit.ID,
			&permit.ParentID,
			&permit.Comment,
			&permit.EmbedSite,
			&permit.TargetSite,
			&permit.TargetLocation,
			&permit.TokenUseNotAfter,
			&permit.SessionTokenUseNotAfter,
		)

		if err != nil {
			return nil, err
		}

		permits = append(permits, permit)
	}

	return permits, nil
}
