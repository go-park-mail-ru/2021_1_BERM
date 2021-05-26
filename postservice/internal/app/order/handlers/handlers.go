package order

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"post/internal/app/models"
	orderUseCase "post/internal/app/order"
	"post/pkg/httputils"
	"strconv"
)

const (
	ctxKeyReqID    uint8 = 1
	ctxUserID      uint8 = 2
	ctxExecutor    uint8 = 3
	ctxQueryParams uint8 = 4
)

type Handlers struct {
	useCase orderUseCase.UseCase
}

func NewHandler(useCase orderUseCase.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func (h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	id := r.Context().Value(ctxUserID).(uint64)
	o := &models.Order{}
	if err := json.NewDecoder(r.Body).Decode(o); err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	o.CustomerID = id
	o, err := h.useCase.Create(*o, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusCreated, o)
}

func (h *Handlers) GetActualOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	param := make(map[string]interface{})
	param["search_str"] = r.URL.Query().Get("search_str")
	if searchStr := r.URL.Query().Get("search_str"); searchStr != "" {
		param["search_str"] = searchStr
	} else {
		param["search_str"] = "~"
	}
	if budgetFrom := r.URL.Query().Get("from"); budgetFrom != "" {
		budgetFromInt, err := strconv.Atoi(budgetFrom)
		if err == nil {
			param["from"] = budgetFromInt
		}
	} else {
		param["from"] = 0
	}
	if budgetTo := r.URL.Query().Get("to"); budgetTo != "" {
		budgetToInt, err := strconv.Atoi(budgetTo)
		if err == nil {
			param["to"] = budgetToInt
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
	ctx := context.WithValue(r.Context(), ctxQueryParams, param)
	o, err := h.useCase.GetActualOrders(ctx)

	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, o)
}

func (h *Handlers) GetOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	o, err := h.useCase.FindByID(id, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, o)
}

func (h *Handlers) ChangeOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	order := models.Order{}
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	params := mux.Vars(r)
	var err error
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	order, err = h.useCase.ChangeOrder(order, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, order)
}

func (h *Handlers) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	err = h.useCase.DeleteOrder(id, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, r, reqID, http.StatusOK, emptyInterface)
}

func (h *Handlers) SelectExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	order := models.Order{}
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	var err error
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	err = h.useCase.SelectExecutor(order, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, order)
}

func (h *Handlers) DeleteExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	order := models.Order{}

	var err error
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	err = h.useCase.DeleteExecutor(order, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, r, reqID, http.StatusOK, emptyInterface)
}

func (h *Handlers) GetAllUserOrders(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}

	o, err := h.useCase.FindByUserID(userID, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, o)
}

func (h *Handlers) CloseOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	orderID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}

	err = h.useCase.CloseOrder(orderID, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, r, reqID, http.StatusOK, emptyInterface)
}

func (h *Handlers) GetAllArchiveUserOrders(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	userInfo := models.UserBasicInfo{}
	var err error
	userInfo.ID, err = strconv.ParseUint(params["id"], 10, 64)
	userInfo.Executor = r.Context().Value(ctxExecutor).(bool)
	o, err := h.useCase.GetArchiveOrders(userInfo, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, o)
}

func (h *Handlers) SearchOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	orderSearch := models.OrderSearch{}
	if err := json.NewDecoder(r.Body).Decode(&orderSearch); err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	o, err := h.useCase.SearchOrders(orderSearch.Keyword, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, o)
}
