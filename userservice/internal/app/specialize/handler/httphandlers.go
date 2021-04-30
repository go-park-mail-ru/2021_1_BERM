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

func New(specializeUseCase usecase.UseCase) *Handler {
	return &Handler{
		specializeUseCase: specializeUseCase,
	}
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	err = h.specializeUseCase.Remove(id, context.Background())
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	httputils.Respond(w, reqID, 200, nil)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	s := &models.Specialize{}
	if err := json.NewDecoder(r.Body).Decode(s); err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	err = h.specializeUseCase.AssociateWithUser(id, s.Name, context.Background())
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	u, err := h.userUseCase.GetById(id, context.Background())
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	httputils.Respond(w, reqID, 200, u)
}
