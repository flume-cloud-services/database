package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/syndtr/goleveldb/leveldb"
)

type Insert struct {
	Database string `json:"database"`
	Table    string `json:"table"`
	Content  string `json:"content"`
}

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

	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)

	var contentResult []map[string]interface{}
	json.Unmarshal([]byte(content.Content), &contentResult)
	result[content.Table] = contentResult

	dataByte, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Put([]byte(content.Database), dataByte, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Data successfully inserted into table " + content.Table)
}
