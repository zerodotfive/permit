package permits

import (
	"time"
	"uuid"
)

type PermitShort struct {
	ID                      uuid.UUID  `json:"id" db:"id"`
	ParentID                *uuid.UUID `json:"parent_id" db:"parent_id"`
	Comment                 string     `json:"comment" db:"comment"`
	EmbedSite               string     `json:"embed_site" db:"embed_site"`
	TargetSite              string     `json:"target_site" db:"target_site"`
	TargetLocation          string     `json:"target_location" db:"target_location"`
	TokenUseNotAfter        time.Time  `json:"token_use_not_after" db:"token_use_not_after"`
	SessionTokenUseNotAfter *time.Time `json:"session_token_use_not_after" db:"session_token_use_not_after"`
}

type PermitToken struct {
	PermitShort
	Token []byte `json:"token"`
}

type PermitSession struct {
	PermitShort
	ChainToken   []byte `json:"chain_token"`
	SessionToken []byte `json:"session_token"`
}
