package router

import (
	"leadgen/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/generate-lead", controller.GenerateLead).Methods("POST")

	return router
}
