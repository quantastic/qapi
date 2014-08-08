package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/quantastic/qapi"
)

func Router() http.Handler {
	mux := mux.NewRouter()
	mux.Methods("GET").Path("/").HandlerFunc(Home)
	mux.Methods("POST").Path("/times").HandlerFunc(CreateTime)
	return mux
}

func Home(w http.ResponseWriter, r *http.Request) {
	jsonWrite(w, http.StatusOK, map[string]interface{}{})
}

func CreateTime(w http.ResponseWriter, r *http.Request) {
	t := &qapi.Time{}
	if err := jsonDecode(w, r, t); err != nil {
		return
	}
	jsonWrite(w, http.StatusCreated, map[string]interface{}{"time": t})
}
