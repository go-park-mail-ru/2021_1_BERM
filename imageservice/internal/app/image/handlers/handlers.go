package handlers

import (
	"encoding/json"
	"image/internal/app/httputils"
	imgUseCase "image/internal/app/image/usecase"
	"image/internal/app/models"
	"net/http"
)

const (
	ctxKeySession uint8 = iota
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
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	defer r.Body.Close()
	u := models.UserImg{}
	err := json.NewDecoder(r.Body).Decode(&u)
	//TODO: вытащить айдишник
	u.ID = r.Context().Value(ctxKeySession).(*models.Session).UserID
	if err != nil {
		httputils.RespondError(w, reqId, InvalidJSON)
		return
	}

	u, err = h.useCase.SetImage(u)
	if err != nil {
		httputils.RespondError(w, reqId, New(err))
		return
	}
	httputils.Respond(w, reqId, http.StatusOK, u)
}
