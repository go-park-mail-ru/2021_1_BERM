package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"user/internal/app/models"
	"user/internal/app/user/usecase"
	"user/pkg/httputils"
)

const (
	ctxKeyReqID uint8 = 1
)

type Handlers struct {
	userUseCase usecase.UseCase
}

func New(userUseCase usecase.UseCase) *Handlers {
	return &Handlers{
		userUseCase: userUseCase,
	}
}

func (h *Handlers) ChangeProfile(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	u := &models.ChangeUser{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	u.ID = id
	response, err := h.userUseCase.Change(*u, context.Background())
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, response)
}

func (h *Handlers) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	ID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	u, err := h.userUseCase.GetById(ID, context.Background())
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, u)
}
