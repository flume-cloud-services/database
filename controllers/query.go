package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/oliveagle/jsonpath"
	"github.com/syndtr/goleveldb/leveldb"
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

	var json_data interface{}
	json.Unmarshal([]byte(data), &json_data)
	pat, _ := jsonpath.Compile(queryStr.Content)
	res, err := pat.Lookup(json_data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(res)

}
