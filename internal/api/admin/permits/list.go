package permits

import (
	"context"
	"net/http"

	"github.com/zerodotfive/appcore/auth/types"
	"github.com/zerodotfive/permit/internal/permits"
)

type ModelPermitListResponse struct {
	Status int
	Body   *ModelPermitListResponseBody
}

type ModelPermitListResponseBody struct {
	Status  bool                  `json:"status"`
	Permits []permits.PermitShort `json:"permits"`
}
type ModelPermitListRequest struct {
	RemoteIP string `header:"X-Forwarded-For" required:"true"`
}

var respListError = &ModelPermitListResponse{
	Status: http.StatusInternalServerError,
	Body: &ModelPermitListResponseBody{
		Status: false,
	},
}

func (h *Handler) list(ctx context.Context, request *ModelPermitListRequest) (*ModelPermitListResponse, error) {
	userInfo, ok := ctx.Value("userInfo").(*types.UserInfo)
	if !ok {
		h.logger.Printf("remoteIP=%s admin.list: user info not in context", request.RemoteIP)
		return respListError, nil
	}

	provider := permits.NewProvider(ctx, h.postgres)

	p, err := provider.List()
	if err != nil {
		h.logger.Printf("remoteIP=%s user=%s admin.list: error list permits: %v", request.RemoteIP, userInfo.Login, err)
		return respListError, nil
	}

	h.logger.Printf("remoteIP=%s user=%s admin.list: permits list success", request.RemoteIP, userInfo.Login)
	return &ModelPermitListResponse{
		Status: http.StatusOK,
		Body: &ModelPermitListResponseBody{
			Status:  true,
			Permits: p,
		},
	}, nil
}
