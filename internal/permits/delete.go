package permits

import "uuid"

func (p *Provider) Delete(id uuid.UUID) error {
	_, err := p.postgres.Exec(p.ctx, `DELETE FROM permits WHERE id=$1;`, id)

	return err
}
