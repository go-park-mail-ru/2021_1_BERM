package image

import (
	"encoding/json"
	"imageservice/internal/app/httputils"
	imgUseCase "imageservice/internal/app/image"
	"imageservice/internal/app/models"
	"net/http"
)

const (
	ctxKeyReqID   uint8 = 1
)

type Handlers struct {
	useCase imgUseCase.UseCase
}

func NewHandler(useCase imgUseCase.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func (h *Handlers) PutAvatar(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	u := models.UserImg{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	u, err = h.useCase.SetImage(u)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, u)
}
