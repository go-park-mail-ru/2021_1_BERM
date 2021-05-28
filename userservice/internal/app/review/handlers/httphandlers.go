package handlers

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"user/internal/app/models"
	"user/internal/app/review"
	"user/pkg/httputils"
	"user/pkg/types"
)

const (
	ctxKeyReqID types.CtxKey = 1
)

type Handler struct {
	reviewsUseCase review.UseCase
}

func New(reviewsUseCase review.UseCase) *Handler {
	return &Handler{
		reviewsUseCase: reviewsUseCase,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	review := &models.Review{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	if err := review.UnmarshalJSON(body); err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}

	review, err = h.reviewsUseCase.Create(*review, r.Context())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	result, err := review.MarshalJSON()
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, 200, result)
}

func (h *Handler) GetAllByUserId(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	ID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}

	reviews, err := h.reviewsUseCase.GetAllReviewByUserId(ID, r.Context())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	result, err := reviews.MarshalJSON()
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, 200, result)
}
