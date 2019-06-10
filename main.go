package main

import (
	"log"
	"net/http"

	"github.com/flume-cloud-services/database/controllers"
)

func main() {
	http.HandleFunc("/signin", controllers.Signin)

	log.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
