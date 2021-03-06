package main

import (
	"io"
	"log"
	"strings"
	"net/http"

	"github.com/flume-cloud-services/database/controllers"
	"github.com/flume-cloud-services/database/middleware"
)

func main() {
	http.HandleFunc("/signin", controllers.Signin)

	http.Handle("/welcome", middleware.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "Welcome to you visitor !")
		}),
		middleware.AuthMiddleware,
	))

	http.Handle("/database/create", middleware.Middleware(
		http.HandlerFunc(controllers.CreateDatabase),
		middleware.AuthMiddleware,
	))
	http.Handle("/database/delete", middleware.Middleware(
		http.HandlerFunc(controllers.DeleteDatabase),
		middleware.AuthMiddleware,
	))

	http.Handle("/query", middleware.Middleware(
		http.HandlerFunc(controllers.CreateQuery),
		middleware.AuthMiddleware,
	))

	http.Handle("/insert", middleware.Middleware(
		http.HandlerFunc(controllers.InsertData),
		middleware.AuthMiddleware,
	))

	http.HandleFunc("/dash/login", controllers.Login)
	http.HandleFunc("/dash/json", controllers.GetData)

	http.Handle("/dash/", http.StripPrefix(strings.TrimRight("/dash/", "/"), http.FileServer(http.Dir("dash"))))

	log.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
