package vacancy

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"post/internal/app/models"
	vacancyUseCase "post/internal/app/vacancy"
	"post/pkg/httputils"
	"strconv"
)

const (
	ctxKeySession uint8 = iota
	ctxKeyReqID   uint8 = 1
	ctxUserID     uint8 = 2
	ctxParam      uint8 = 4
	ctxExecutor   uint8 = 3

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

	id := r.Context().Value(ctxUserID).(uint64)
	v := &models.Vacancy{
		CustomerID: id,
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	var err error
	if v, err = h.useCase.Create(*v, context.Background()); err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusCreated, v)
}

func (h *Handlers) GetVacancy(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	v, err := h.useCase.FindByID(id, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, v)
}

func (h *Handlers) GetActualVacancies(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	param := make(map[string]interface{})
	param["search_str"] = r.URL.Query().Get("search_str")
	if searchStr := r.URL.Query().Get("search_str"); searchStr != "" {
		param["search_str"] = searchStr
	} else {
		param["search_str"] = "~"
	}
	if salaryFrom := r.URL.Query().Get("from"); salaryFrom != "" {
		salaryFromInt, err := strconv.Atoi(salaryFrom)
		if err == nil {
			param["from"] = salaryFromInt
		}
	} else {
		param["from"] = 0
	}
	if salaryTo := r.URL.Query().Get("to"); salaryTo != "" {
		salaryToInt, err := strconv.Atoi(salaryTo)
		if err == nil {
			param["to"] = salaryToInt
		}
	} else {
		param["to"] = 0
	}

	if desc := r.URL.Query().Get("desc"); desc != "" {
		descBool, err := strconv.ParseBool(desc)
		if err == nil {
			param["desc"] = descBool
		}
	} else {
		param["desc"] = false
	}

	if category := r.URL.Query().Get("category"); category != "" {
		param["category"] = category
	} else {
		param["category"] = "~"
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err == nil {
			param["limit"] = limitInt
		}
	} else {
		param["limit"] = 15
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err == nil {
			param["offset"] = offsetInt
		}
	} else {
		param["offset"] = 0
	}

	v, err := h.useCase.GetActualVacancies(context.WithValue(r.Context(), ctxParam, param))
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	var vReq map[string]interface{}
	vReq["vacancies"] = v;
	vReq["size"] = len(v);
	httputils.Respond(w, r, reqID, http.StatusOK, vReq)
}

func (h *Handlers) ChangeVacancy(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	vacancy := models.Vacancy{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&vacancy); err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	params := mux.Vars(r)
	vacancy.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	vacancy, err = h.useCase.ChangeVacancy(vacancy, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, vacancy)
}

func (h *Handlers) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	err = h.useCase.DeleteVacancy(id, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, r, reqID, http.StatusOK, emptyInterface)
}

func (h *Handlers) GetAllUserVacancies(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	vacancies, err := h.useCase.FindByUserID(userID, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, vacancies)
}

func (h *Handlers) SelectExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	vacancy := models.Vacancy{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&vacancy); err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	vacancy.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	err = h.useCase.SelectExecutor(vacancy, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, vacancy)
}

func (h *Handlers) DeleteExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	vacancy := models.Vacancy{}

	var err error
	vacancy.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	err = h.useCase.DeleteExecutor(vacancy, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, r, reqID, http.StatusOK, emptyInterface)
}

func (h *Handlers) CloseVacancy(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	vacancyID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}

	err = h.useCase.CloseVacancy(vacancyID, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, r, reqID, http.StatusOK, emptyInterface)
}

func (h *Handlers) GetAllArchiveUserVacancies(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userInfo := models.UserBasicInfo{}
	var err error
	userInfo.ID, err = strconv.ParseUint(params["id"], 10, 64)
	userInfo.Executor = r.Context().Value(ctxExecutor).(bool)

	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	v, err := h.useCase.GetArchiveVacancies(userInfo, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, v)
}

func (h *Handlers) SearchVacancy(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	vacancySearch := models.VacancySearch{}
	if err := json.NewDecoder(r.Body).Decode(&vacancySearch); err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	v, err := h.useCase.SearchVacancy(vacancySearch.Keyword, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, v)
}

func (h *Handlers) SuggestVacancyTitle(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	suggestWord := r.URL.Query().Get("suggest_word")
	suggestTitles, err := h.useCase.SuggestVacancyTitle(suggestWord, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, suggestTitles)
}
