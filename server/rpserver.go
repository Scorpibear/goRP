package server

import (
	"github.com/reportportal/goRP/conf"
	"github.com/reportportal/goRP/registry"
	"goji.io"
	"goji.io/pat"

	"github.com/reportportal/goRP/commons"
	"log"
	"net/http"
	"strconv"
)

//RpServer represents ReportPortal micro-service instance
type RpServer struct {
	mux *goji.Mux
	cfg *conf.RpConfig
	Sd  registry.ServiceDiscovery
}

//New creates new instance of RpServer struct
func New(cfg *conf.RpConfig) *RpServer {

	log.Println(commons.Build)

	var sd registry.ServiceDiscovery
	switch cfg.Registry {
	case conf.Eureka:
		sd = registry.NewEureka(cfg)
	case conf.Consul:
		cfg.Consul.Tags = cfg.Consul.Tags + ",statusPageUrlPath=/info" + "," + "healthCheckUrlPath=/health"
		sd = registry.NewConsul(cfg)
	}

	srv := &RpServer{
		mux: goji.NewMux(),
		cfg: cfg,
		Sd:  sd,
	}

	srv.mux.HandleFunc(pat.Get("/health"), func(w http.ResponseWriter, rq *http.Request) {
		commons.WriteJSON(200, map[string]string{"status": "UP"}, w)
	})

	commons.Build.Name = cfg.AppName
	srv.mux.HandleFunc(pat.Get("/info"), func(w http.ResponseWriter, rq *http.Request) {
		commons.WriteJSON(200, commons.Build, w)

	})
	return srv
}

//AddRoute gives access to GIN router to add route and perform other modifications
func (srv *RpServer) AddRoute(f func(router *goji.Mux)) {
	f(srv.mux)
}

//StartServer starts HTTP server
func (srv *RpServer) StartServer() {

	if nil != srv.Sd {
		registry.Register(srv.Sd)
	}
	// listen and server on mentioned port
	log.Printf("Starting on port %d", srv.cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(srv.cfg.Server.Port), srv.mux))
}
