package handler

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"user/internal/app/specialize/usecase"
	"user/pkg/httputils"
)

const (
	ctxKeyReqID uint8 = 1
)

type Handler struct {
	specializeUseCase usecase.UseCase
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

func (h *Handler) GetSpecialize(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	spec, err := h.specializeUseCase.FindByUseID(id, context.Background())
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	httputils.Respond(w, reqID, 200, map[string][]string{
		"specializes": spec,
	})
}
