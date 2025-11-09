package server

import (
	"encoding/json"
	"github.com/G4GANCHAUDHARY/encryption-service/global"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/core"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/dataMapper/httpDataMapper"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/httpModel"
	"github.com/G4GANCHAUDHARY/encryption-service/utils"
	mux "github.com/gorilla/mux"
	"github.com/lib/pq"
	"net/http"
	"strconv"
)

type HttpServer struct {
	Router                *mux.Router
	Core                  core.IUrlShortenerCore
	HttpRequestDataMapper httpDataMapper.IHttpRequestDataMapper
}

func (h HttpServer) Init() error {
	h.registerRoutes()
	return nil
}

func (h HttpServer) registerRoutes() {
	h.Router.HandleFunc("/url-shortener/v1", h.UrlHandler).Methods(http.MethodPost, http.MethodGet)
	h.Router.HandleFunc("/{short_code}", h.HandleGetUrlDetails).Methods(http.MethodGet)
}

func (h *HttpServer) UrlHandler(rw http.ResponseWriter, rq *http.Request) {
	ctx := utils.CreateServerContext(rq)
	defer utils.HandlePanic(utils.HandlePanicRequest{
		FuncName: "UrlHandler",
		CallbackFn: func() {
			utils.HTTPFailWithXxx(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, rw)
		},
	})

	switch rq.Method {
	case http.MethodPost:
		var httpPayload httpModel.GenerateUrlReqPayload
		err := json.NewDecoder(rq.Body).Decode(&httpPayload)
		if err != nil {
			utils.HTTPFailWithXxx("bad request", http.StatusBadRequest, rw)
			return
		}

		if err = httpPayload.Validate(); err != nil {
			utils.HTTPFailWithXxx(err.Error(), http.StatusBadRequest, rw)
			return
		}

		res, err := h.Core.EncryptUrl(ctx, h.HttpRequestDataMapper.GetGenerateUrlCoreReq(httpPayload))
		if err != nil {
			utils.HTTPFailWithXxx(err.Error(), http.StatusInternalServerError, rw)
			return
		}
		utils.HTTPSuccessWith200(res, rw)

	case http.MethodGet:
		pageSize, err := strconv.Atoi(rq.URL.Query().Get("page_size"))
		if err != nil {
			pageSize = global.DefaultPageSize
		}
		cursor, err := strconv.Atoi(rq.URL.Query().Get("cursor"))
		if err != nil {
			cursor = global.DefaultCursor
		}
		orderBy, sortingOrder, err := utils.GetSortParams(rq)
		filters := pq.StringArray{"url.is_active = false"}

		coreRequest := h.HttpRequestDataMapper.GetUrlListCoreReq(pageSize, cursor, sortingOrder, filters, orderBy)
		res, err := h.Core.GetUrls(ctx, coreRequest)
		if err != nil {
			utils.HTTPFailWithXxx(err.Error(), http.StatusInternalServerError, rw)
			return
		}
		utils.HTTPSuccessWith200(res, rw)
	}

	return
}

func (h *HttpServer) HandleGetUrlDetails(rw http.ResponseWriter, rq *http.Request) {
	ctx := utils.CreateServerContext(rq)

	defer utils.HandlePanic(utils.HandlePanicRequest{
		FuncName: "HandleGetUrlDetails",
		CallbackFn: func() {
			utils.HTTPFailWithXxx(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, rw)
		},
	})

	switch rq.Method {

	case http.MethodGet:
		shortCode, found := mux.Vars(rq)["short_code"]
		if found == false {
			utils.HTTPFailWithXxx("short-url not present", http.StatusBadRequest, rw)
			return
		}

		res, err := h.Core.DecryptUrl(ctx, h.HttpRequestDataMapper.GetDecryptUrlCoreReq(shortCode))
		if err != nil {
			utils.HTTPFailWithXxx(err.Error(), http.StatusNotFound, rw)
			return
		}
		utils.HTTPSuccessWith200(res, rw)
	}
}
