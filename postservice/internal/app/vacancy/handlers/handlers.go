package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"post/internal/app/models"
	vacancyUseCase "post/internal/app/vacancy/usecase"
	"post/pkg/httputils"
	"strconv"
)

const (
	ctxKeySession uint8 = iota
	ctxKeyReqID   uint8 = 1
	ctxUserInfo   uint8 = 2
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
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err)
			return
		}
	}(r.Body)
	id := r.Context().Value(ctxUserInfo).(uint64)
	v := &models.Vacancy{
		CustomerID: id,
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	var err error
	if v, err = h.useCase.Create(*v); err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusCreated, v)
}

func (h *Handlers) GetVacancy(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err)
			return
		}
	}(r.Body)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	v, err := h.useCase.FindByID(id)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, v)
}

func (h *Handlers) GetActualVacancies(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	v, err := h.useCase.GetActualVacancies()
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, v)
}

func (h *Handlers) ChangeVacancy(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err)
			return
		}
	}(r.Body)
	vacancy := models.Vacancy{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&vacancy); err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	params := mux.Vars(r)
	vacancy.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	vacancy, err = h.useCase.ChangeVacancy(vacancy)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, vacancy)
}

func (h *Handlers) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err)
			return
		}
	}(r.Body)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	err = h.useCase.DeleteVacancy(id)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, reqID, http.StatusOK, emptyInterface)
}

func (h *Handlers) GetAllUserVacancies(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err)
			return
		}
	}(r.Body)
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	vacancies, err := h.useCase.FindByUserID(userID)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, vacancies)
}

func (h *Handlers) SelectExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			httputils.RespondError(w, reqID, err)
			return
		}
	}(r.Body)
	vacancy := models.Vacancy{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&vacancy); err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	vacancy.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	err = h.useCase.SelectExecutor(vacancy)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, vacancy)
}

func (h *Handlers) DeleteExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	vacancy := models.Vacancy{}

	var err error
	vacancy.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	err = h.useCase.DeleteExecutor(vacancy)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, reqID, http.StatusOK, emptyInterface)
}
