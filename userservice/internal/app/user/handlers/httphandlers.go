package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"user/internal/app/models"
	"user/internal/app/user"
	"user/pkg/httputils"
)

const (
	ctxKeyReqID uint8 = 1
	ctxParam uint8 = 4
)

type Handlers struct {
	userUseCase user.UseCase
}

func New(userUseCase user.UseCase) *Handlers {
	return &Handlers{
		userUseCase: userUseCase,
	}
}

func (h *Handlers) ChangeProfile(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	u := &models.ChangeUser{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	u.ID = id
	response, err := h.userUseCase.Change(*u, r.Context())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, response)
}

func (h *Handlers) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	ID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	u, err := h.userUseCase.GetById(ID, r.Context())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, u)
}

func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	param := make(map[string]interface{})
	param["search_str"] = r.URL.Query().Get("search_str")
	if searchStr := r.URL.Query().Get("search_str"); searchStr != "" {
		param["search_str"] = searchStr
	} else {
		param["search_str"] = "~"
	}

	if searchStr := r.URL.Query().Get("search_str"); searchStr != "" {
		param["search_str"] = searchStr
	} else {
		param["search_str"] = "~"
	}
	if sort := r.URL.Query().Get("sort"); sort != ""{
		param["sort"] = sort;
	}else{
		param[sort] = "~";
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
		param["category"] = ""
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
	u, err := h.userUseCase.GetUsers(context.WithValue(r.Context(), ctxParam, param))
	if err != nil{
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, u)
}
