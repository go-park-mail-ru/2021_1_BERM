package delivery

import (
	"encoding/json"
	"ff/internal/app/models"
	"ff/internal/app/vacancy"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ctxKey uint8

const (
	ctxKeySession ctxKey = iota
	ctxKeyReqID   ctxKey = 1
)

type VcancyHandler struct {
	vacancyUseCase vacancy.VacancyUseCase
}
func (vac *VcancyHandler) handleCreateVacancy(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	id := r.Context().Value(ctxKeySession).(*models.Session).UserID
	v := &models.Vacancy{
		UserID: id,
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		vac.error(w, reqId, InvalidJSON)
		return
	}
	var err error
	if v, err = vac.vacancyUseCase.Create(*v); err != nil {
		vac.error(w, reqId, New(err))
	}
	respond(w, reqId, http.StatusCreated, v)
}

func (vac *VcancyHandler) handleGetVacancy(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		vac.error(w, reqId, InvalidJSON)
		return
	}
	v, err := vac.vacancyUseCase.FindByID(id)
	if err != nil {
		vac.error(w, reqId, New(err))
		return
	}
	vac.respond(w, reqId, http.StatusOK, v)
}

func (vac *VcancyHandler) handleCreateVacancyResponse(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	response := &models.ResponseVacancy{}
	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		vac.error(w, http.StatusBadRequest, InvalidJSON) //Bad json
		return
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		vac.error(w, http.StatusBadRequest, New(err)) //Bad json
		return
	}
	response.VacancyID = id
	response, err = vac.vacancyUseCase.Create(*response)
	if err != nil {
		vac.error(w, http.StatusBadRequest, New(err))
		return
	}
	vac.respond(w, reqId, http.StatusCreated, response)
}

func (vac *VcancyHandler) handleGetAllVacancyResponses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		vac.error(w, http.StatusBadRequest, InvalidJSON) //Bad json
		return
	}
	responses, err := vac.vacancyUseCase.FindByVacancyID(id)
	if err != nil {
		vac.error(w, http.StatusBadRequest, New(err)) //Bad json
		return
	}

	vac.respond(w, reqId, http.StatusOK, responses)
}
