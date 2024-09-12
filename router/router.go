package router

import (
	"leadgen/controller"

	"github.com/gin-gonic/gin"
)

// func Router() *mux.Router {

// 	router := mux.NewRouter()

// 	router.HandleFunc("/generate-lead", controller.GenerateLead).Methods("POST")
// 	router.HandleFunc("/upload-csv", controller.UploadCSV).Methods("POST")

// 	return router
// }

func Router() *gin.Engine {
	router := gin.Default()

	//Route for generating a Lead
	router.POST("/generate-lead", func(c *gin.Context) {
		controller.GenerateLead(c)
	})

	router.POST("/upload-csv", func(c *gin.Context) {
		controller.UploadCSV(c)
	})

	return router
}
