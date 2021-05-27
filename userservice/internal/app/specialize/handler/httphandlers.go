package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"user/internal/app/models"
	"user/internal/app/specialize"
	"user/internal/app/user"
	"user/pkg/httputils"
	"user/pkg/types"
)

const (
	ctxKeyReqID types.CtxKey = 1
)

type Handler struct {
	specializeUseCase specialize.UseCase
	userUseCase       user.UseCase
}

func New(specializeUseCase specialize.UseCase, userUseCase user.UseCase) *Handler {
	return &Handler{
		specializeUseCase: specializeUseCase,
		userUseCase:       userUseCase,
	}
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	s := &models.Specialize{}
	if err := json.NewDecoder(r.Body).Decode(s); err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	err = h.specializeUseCase.Remove(id, s.Name, r.Context())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, 200, nil)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	s := &models.Specialize{}
	if err = json.NewDecoder(r.Body).Decode(s); err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	err = h.specializeUseCase.AssociateWithUser(id, s.Name, r.Context())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	u, err := h.userUseCase.GetById(id, r.Context())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, 200, u)
}
