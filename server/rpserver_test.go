package server

import (
	"github.com/gorilla/handlers"
	"github.com/reportportal/goRP/commons"
	"github.com/reportportal/goRP/conf"
	"goji.io"
	"goji.io/pat"
	"net/http"
	"os"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func ExampleRpServer() {
	rpConf := conf.LoadConfig("../server.yaml", nil)
	rp := New(rpConf)

	rp.AddRoute(func(router *goji.Mux) {
		router.HandleFunc(pat.Get("/ping"), func(w http.ResponseWriter, rq *http.Request) {
			commons.WriteJSON(http.StatusOK, Person{"av", 20}, w)
		})
	})

	rp.StartServer()

}

func ExampleRpServer_StartServer() {

	rpConf := conf.LoadConfig("../server.yaml",
		map[string]interface{}{"AuthServerURL": "http://localhost:9998/sso/me"})

	srv := New(rpConf)

	srv.AddRoute(func(mux *goji.Mux) {
		mux.Use(func(next http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, next)
		})

		secured := goji.SubMux()
		secured.Use(RequireRole("USER", rpConf.Get("AuthServerURL").(string)))

		me := func(w http.ResponseWriter, rq *http.Request) {
			commons.WriteJSON(http.StatusOK, rq.Context().Value("user"), w)

		}
		secured.HandleFunc(pat.Get("/me"), me)

		mux.Handle(pat.Get("/"), secured)

	})

	srv.StartServer()
}
