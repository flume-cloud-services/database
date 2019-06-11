package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tidwall/sjson"
)

// Insert struct
type Insert struct {
	Database string `json:"database"`
	Query    string `json:"query"`
	Content  string `json:"content"`
}

// InsertData into a specified table into a specified database
func InsertData(w http.ResponseWriter, r *http.Request) {
	var content Insert
	err := json.NewDecoder(r.Body).Decode(&content)
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

	data, err := db.Get([]byte(content.Database), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := sjson.Set(string(data), content.Query, content.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Put([]byte(content.Database), []byte(res), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Data successfully inserted into table ")
}
