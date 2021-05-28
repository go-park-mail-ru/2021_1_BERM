package response

import (
	"context"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"post/internal/app/models"
	responseUseCase "post/internal/app/response"
	"post/pkg/httputils"
	"post/pkg/types"
	"strconv"
)

const (
	ctxKeyReqID types.CtxKey = 1
	ctxUserID   types.CtxKey = 2
)

type Handlers struct {
	useCase responseUseCase.UseCase
}

func NewHandler(useCase responseUseCase.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func (h *Handlers) CreatePostResponse(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	response := &models.Response{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	if err := response.UnmarshalJSON(body); err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	response.PostID = id
	response.VacancyResponse = r.URL.String() == "/api/vacancy/"+strconv.FormatUint(id, 10)+"/response"
	response.OrderResponse = r.URL.String() == "/api/order/"+strconv.FormatUint(id, 10)+"/response"
	response, err = h.useCase.Create(*response, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	result, err := response.MarshalJSON()
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, http.StatusCreated, result)
}

func (h *Handlers) GetAllPostResponses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	vacancyResponse := r.URL.String() == "/api/vacancy/"+strconv.FormatUint(id, 10)+"/response"
	orderResponse := r.URL.String() == "/api/order/"+strconv.FormatUint(id, 10)+"/response"
	responses := models.ResponseList{}
	responses, err = h.useCase.FindByPostID(id, orderResponse, vacancyResponse, context.Background())
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}

	result, err := responses.MarshalJSON()
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, result)
}

func (h *Handlers) ChangePostResponse(w http.ResponseWriter, r *http.Request) {
	response := &models.Response{}
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	if err := response.UnmarshalJSON(body); err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}

	params := mux.Vars(r)
	response.PostID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}
	response.UserID = r.Context().Value(ctxUserID).(uint64)
	response.VacancyResponse = r.URL.String() == "/api/vacancy/"+strconv.FormatUint(response.PostID, 10)+"/response"
	response.OrderResponse = r.URL.String() == "/api/order/"+strconv.FormatUint(response.PostID, 10)+"/response"
	responses, err := h.useCase.Change(*response, context.Background())

	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}

	result, err := responses.MarshalJSON()
	if err != nil {
		httputils.RespondError(w, r, reqID, err)
		return
	}
	httputils.Respond(w, r, reqID, http.StatusOK, result)
}

func (h *Handlers) DelPostResponse(w http.ResponseWriter, r *http.Request) {
	response := &models.Response{}
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	var err error
	response.PostID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}

	response.UserID = r.Context().Value(ctxUserID).(uint64)
	response.VacancyResponse = r.URL.String() == "/api/vacancy/"+strconv.FormatUint(response.PostID, 10)+"/response"
	response.OrderResponse = r.URL.String() == "/api/order/"+strconv.FormatUint(response.PostID, 10)+"/response"
	err = h.useCase.Delete(*response, context.Background())

	if err != nil {
		httputils.RespondError(w, r, reqID, err)

		return
	}

	httputils.Respond(w, r, reqID, http.StatusOK, nil)
}
