package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"user/internal/app/models"
	"user/internal/app/review/usecase"
	"user/pkg/httputils"
)

const (
	ctxKeyReqID uint8 = 1
)

type Handler struct {
	reviewsUseCase usecase.UseCase
}

func New(reviewsUseCase usecase.UseCase) *Handler {
	return &Handler{
		reviewsUseCase: reviewsUseCase,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	review := &models.Review{}
	if err := json.NewDecoder(r.Body).Decode(review); err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	review, err := h.reviewsUseCase.Create(*review, r.Context())
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, 200, review)
}


func (h *Handler) GetAllByUserId(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	ID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}

	reviews, err := h.reviewsUseCase.GetAllReviewByUserId(ID, r.Context())
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, 200, reviews)
}