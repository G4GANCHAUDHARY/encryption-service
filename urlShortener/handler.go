package urlShortener

import (
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/core"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/dataMapper/httpDataMapper"
	server2 "github.com/G4GANCHAUDHARY/encryption-service/urlShortener/server"
	"github.com/gorilla/mux"
)

type UrlHandler struct {
	Router                *mux.Router
	Core                  core.IUrlShortenerCore
	HttpRequestDataMapper httpDataMapper.IHttpRequestDataMapper
}

func (uh *UrlHandler) Init() {
	server := server2.HttpServer{
		Router:                uh.Router,
		Core:                  uh.Core,
		HttpRequestDataMapper: uh.HttpRequestDataMapper,
	}
	server.Init()
}
