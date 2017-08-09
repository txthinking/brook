package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// RunHTTPServer will start a http server
func RunHTTPServer(address string) error {
	r := mux.NewRouter()

	r.Methods("POST").Path("/setting").HandlerFunc(PostSetting)
	r.Methods("GET").Path("/setting").HandlerFunc(GetSetting)
	r.Methods("GET").Path("/pac").HandlerFunc(GetPAC)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewStatic(assetFS()))
	n.UseHandler(r)

	return http.ListenAndServe(address, n)
}
