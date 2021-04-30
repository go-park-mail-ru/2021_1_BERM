package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"post/internal/app/models"
	vacancyUseCase "post/internal/app/vacancy/usecase"
	"post/pkg/Error"
	"post/pkg/httputils"
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

func (h *Handlers) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		}
	}(r.Body)
	id, err := strconv.ParseUint(r.Header.Get("X_Id"), 10, 64)
	v := &models.Vacancy{
		CustomerID: id,
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

func (h *Handlers) GetVacancy(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		}
	}(r.Body)
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

func (h *Handlers) ChangeVacancy(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		}
	}(r.Body)
	vacancy := models.Vacancy{}
	if err = json.NewDecoder(r.Body).Decode(&vacancy); err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	vacancy.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	vacancy, err = h.useCase.ChangeVacancy(vacancy)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, vacancy)
}

func (h *Handlers) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		}
	}(r.Body)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	err = h.useCase.DeleteVacancy(id)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, reqID, http.StatusOK, emptyInterface)
}

func (h *Handlers) GetAllUserVacancies(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		}
	}(r.Body)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	vacancies, err := h.useCase.FindByUserID(userID)
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
	httputils.Respond(w, reqID, http.StatusOK, vacancies)
}

func (h *Handlers) SelectExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		}
	}(r.Body)
	vacancy := models.Vacancy{}
	if err = json.NewDecoder(r.Body).Decode(&vacancy); err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	vacancy.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	err = h.useCase.SelectExecutor(vacancy)
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
	httputils.Respond(w, reqID, http.StatusOK, vacancy)
}

func (h *Handlers) DeleteExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}

	vacancy := models.Vacancy{}

	vacancy.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	err = h.useCase.DeleteExecutor(vacancy)
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
	var emptyInterface interface{}
	httputils.Respond(w, reqID, http.StatusOK, emptyInterface)
}
