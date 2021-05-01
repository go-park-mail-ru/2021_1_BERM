package handlers

import (
	"encoding/json"
	"image/internal/app/httputils"
	imgUseCase "image/internal/app/image/usecase"
	"image/internal/app/models"
	"net/http"
	"strconv"
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
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	u := models.UserImg{}
	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	u.ID, err = strconv.ParseUint(r.Header.Get("X_Id"), 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	u, err = h.useCase.SetImage(u)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, u)
}
