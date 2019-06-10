package controllers

import (
	"net/http"
	"encoding/json"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello")
}