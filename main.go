package main

import (
	"fmt"
	"leadgen/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	log.Fatal(http.ListenAndServe(":6000", r))
	fmt.Println("Listening at port 6000")
}
