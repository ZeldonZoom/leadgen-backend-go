package router

import (
	"leadgen/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/generate-lead", controller.GenerateLead).Methods("POST")
	router.HandleFunc("/upload-csv", controller.UploadCSV).Methods("POST")

	return router
}
