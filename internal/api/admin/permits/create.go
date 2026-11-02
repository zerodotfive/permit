package permits

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"uuid"

	"github.com/zerodotfive/appcore/auth/oidc/types"
	"github.com/zerodotfive/permit/internal/permits"
)

type ModelPermitCreateResponse struct {
	Status int
	Body   *ModelPermitCreateResponseBody
}

type ModelPermitCreateResponseBody struct {
	Status           bool      `json:"status"`
	ID               uuid.UUID `json:"id"`
	URL              string    `json:"url"`
	TokenUseNotAfter time.Time `json:"token_use_not_after"`
}
type ModelPermitCreateRequest struct {
	Body struct {
		Comment                string `json:"comment" required:"true"`
		EmbedSite              string `json:"embed_site" required:"true"`
		TargetSite             string `json:"target_site" required:"true"`
		TargetLocation         string `json:"target_location" required:"false" default:""`
		TokenLifetimeMinutes   int    `json:"token_lifetime_minutes" required:"true" minimum:"1" maximum:"1440"`
		SessionLifetimeMinutes int    `json:"session_lifetime_minutes" required:"true" minimum:"1" maximum:"1440"`
	}
	RemoteIP string `header:"X-Forwarded-For" required:"true"`
}

var respCreateError = &ModelPermitCreateResponse{
	Status: http.StatusInternalServerError,
	Body: &ModelPermitCreateResponseBody{
		Status: false,
	},
}

func (h *Handler) create(ctx context.Context, request *ModelPermitCreateRequest) (*ModelPermitCreateResponse, error) {
	userInfo, ok := ctx.Value("userInfo").(*types.UserInfo)
	if !ok {
		h.logger.Printf("remoteIP=%s admin.create: user info not in context", request.RemoteIP)
		return respCreateError, nil
	}

	esu, err := url.Parse(request.Body.EmbedSite)
	if err != nil || (esu.Scheme != "http" && esu.Scheme != "https") {
		e := fmt.Errorf("remoteIP=%s user=%s admin.create: embed site %s is not a valid url, err: %v", request.RemoteIP, userInfo.Login, request.Body.EmbedSite, err)
		h.logger.Println(e)
		return respCreateError, nil
	}

	tsu, err := url.Parse(request.Body.TargetSite)
	if err != nil || (tsu.Scheme != "http" && tsu.Scheme != "https") {
		e := fmt.Errorf("remoteIP=%s user=%s admin.create: target site %s is not a valid url, err: %v", request.RemoteIP, userInfo.Login, request.Body.TargetSite, err)
		h.logger.Println(e)
		return respCreateError, nil
	}

	provider := permits.NewProvider(ctx, h.postgres)

	permit, err := provider.Create(
		request.Body.Comment,
		request.Body.EmbedSite,
		request.Body.TargetSite,
		request.Body.TargetLocation,
		request.Body.TokenLifetimeMinutes,
		request.Body.SessionLifetimeMinutes,
	)
	if err != nil {
		h.logger.Printf("remoteIP=%s user=%s admin.create: error creating permit: %v", request.RemoteIP, userInfo.Login, err)
		return respCreateError, errors.New("error creating permit")
	}

	u := strings.TrimRight(request.Body.EmbedSite, "/") +
		"/" + strings.Trim(h.locationActivate, "/") + "/" +
		base64.RawURLEncoding.EncodeToString(permit.Token)

	h.logger.Printf("remoteIP=%s user=%s admin.create: permit create success: %s", request.RemoteIP, userInfo.Login, permit.ID.String())
	return &ModelPermitCreateResponse{
		Status: http.StatusOK,
		Body: &ModelPermitCreateResponseBody{
			Status:           true,
			ID:               permit.ID,
			URL:              u,
			TokenUseNotAfter: permit.TokenUseNotAfter,
		},
	}, nil
}
