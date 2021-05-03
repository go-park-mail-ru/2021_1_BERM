package handlers

import (
	"encoding/json"
	"image/internal/app/httputils"
	imgUseCase "image/internal/app/image/usecase"
	"image/internal/app/models"
	"net/http"
)

const (
	ctxKeySession uint8 = 3
	ctxKeyReqID   uint8 = 1
	ctxUserInfo   uint8 = 2
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
	defer r.Body.Close()
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	u := models.UserImg{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	u, err = h.useCase.SetImage(u)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, u)
}
