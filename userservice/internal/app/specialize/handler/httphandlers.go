package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"user/internal/app/models"
	"user/internal/app/specialize/usecase"
	usecase2 "user/internal/app/user/usecase"
	"user/pkg/httputils"
)

const (
	ctxKeyReqID uint8 = 1
)

type Handler struct {
	specializeUseCase usecase.UseCase
	userUseCase usecase2.UseCase
}

func New(specializeUseCase usecase.UseCase, userUseCase usecase2.UseCase) *Handler {
	return &Handler{
		specializeUseCase: specializeUseCase,
		userUseCase: userUseCase,
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
	err = h.specializeUseCase.Remove(id, s.Name, context.Background())
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
	err = h.specializeUseCase.AssociateWithUser(id, s.Name, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	u, err := h.userUseCase.GetById(id, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, 200, u)
}
