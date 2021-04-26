package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"post/internal/app/models"
	vacancyUseCase "post/internal/app/vacancy/usecase"
	"post/pkg/httputils"
	"post/pkg/Error"
	"strconv"
)


const (
	ctxKeySession uint8 = iota
	ctxKeyReqID   uint8 = 1
)

type Handlers struct {
	useCase vacancyUseCase.UseCase
}

func NewHandler(useCase vacancyUseCase.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func(h *Handlers) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}

	id, err  := strconv.ParseUint(r.Header.Get("X_Id"), 10, 64)
	v := &models.Vacancy{
		UserID: id,
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	if v, err = h.useCase.Create(*v); err != nil {

		httpErr := &Error.Error{}
		errors.As(err, &httpErr)
		if httpErr.InternalError {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		} else {
			httputils.RespondError(w, reqID, err, http.StatusBadRequest)
		}
		return
	}
	httputils.Respond(w, reqID, http.StatusCreated, v)
}

func(h *Handlers) GetVacancy(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	v, err := h.useCase.FindByID(id)
	if err != nil {
		httpErr := &Error.Error{}
		errors.As(err, &httpErr)
		if httpErr.InternalError {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		} else {
			httputils.RespondError(w, reqID, err, http.StatusBadRequest)
		}
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, v)
}
