package permits

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (p *Provider) Cleanup(retention time.Duration) (int64, error) {
	tag, err := p.postgres.Exec(p.ctx, `
		DELETE FROM permits
		WHERE (session_token_sha IS NULL AND token_use_not_after < now() - make_interval(secs => $1)) OR
				(session_token_sha IS NOT NULL AND session_token_use_not_after < now() - make_interval(secs => $1));`,
		retention.Seconds(),
	)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func RunCleanup(postgres *pgxpool.Pool, logger *log.Logger, interval, retention time.Duration) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	run := func() {
		runCtx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()

		n, err := NewProvider(runCtx, postgres).Cleanup(retention)
		if err != nil {
			logger.Printf("permits cleanup error: %s", err)
			return
		}
		if n > 0 {
			logger.Printf("permits cleanup: deleted %d", n)
		}
	}

	run()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}
