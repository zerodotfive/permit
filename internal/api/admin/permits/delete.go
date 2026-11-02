package permits

import (
	"context"
	"net/http"
	"uuid"

	"github.com/zerodotfive/appcore/auth/oidc/types"
	"github.com/zerodotfive/permit/internal/permits"
)

type ModelPermitDeleteResponse struct {
	Status int
	Body   *ModelPermitDeleteResponseBody
}

type ModelPermitDeleteResponseBody struct {
	Status bool `json:"status"`
}
type ModelPermitDeleteRequest struct {
	ID       uuid.UUID `path:"id"`
	RemoteIP string    `header:"X-Forwarded-For" required:"true"`
}

var respDeleteError = &ModelPermitDeleteResponse{
	Status: http.StatusInternalServerError,
	Body: &ModelPermitDeleteResponseBody{
		Status: false,
	},
}

func (h *Handler) delete(ctx context.Context, request *ModelPermitDeleteRequest) (*ModelPermitDeleteResponse, error) {
	userInfo, ok := ctx.Value("userInfo").(*types.UserInfo)
	if !ok {
		h.logger.Printf("remoteIP=%s admin.delete: user info not in context", request.RemoteIP)
		return respDeleteError, nil
	}

	provider := permits.NewProvider(ctx, h.postgres)

	err := provider.Delete(request.ID)
	if err != nil {
		h.logger.Printf("remoteIP=%s user=%s admin.delete: error delete permits: %v", request.RemoteIP, userInfo.Login, err)
		return respDeleteError, nil
	}

	h.logger.Printf("remoteIP=%s user=%s admin.delete: permit delete success: %s", request.RemoteIP, userInfo.Login, request.ID.String())
	return &ModelPermitDeleteResponse{
		Status: http.StatusOK,
		Body: &ModelPermitDeleteResponseBody{
			Status: true,
		},
	}, nil
}
