package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tidwall/gjson"
)

// Query struct
type Query struct {
	Name    string `json:"database"`
	Content string `json:"query"`
}

// CreateQuery compile query and lookup for results on a specified database
func CreateQuery(w http.ResponseWriter, r *http.Request) {
	var queryStr Query
	err := json.NewDecoder(r.Body).Decode(&queryStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := leveldb.OpenFile("level.db", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	data, err := db.Get([]byte(queryStr.Name), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := gjson.Get(string(data), queryStr.Content)

	json.NewEncoder(w).Encode(res)

}
