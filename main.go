package main

import (
	"io"
	"log"
	"net/http"

	"github.com/flume-cloud-services/database/controllers"
	"github.com/flume-cloud-services/database/middleware"
)

func main() {
	http.Handle("/welcome", middleware.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "Welcome to you visitor !")
		}),
		middleware.AuthMiddleware,
	))
	http.HandleFunc("/signin", controllers.Signin)

	log.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
