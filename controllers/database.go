package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/syndtr/goleveldb/leveldb"
)

// Database struct to read database name from the request body
type Database struct {
	Name string `json:"name"`
}

// CreateDatabase with the given name
func CreateDatabase(w http.ResponseWriter, r *http.Request) {
	var database Database
	err := json.NewDecoder(r.Body).Decode(&database)
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

	err = db.Put([]byte(database.Name), []byte("{}"), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Database " + database.Name + " succesfully created")
}

// DeleteDatabase with the given name
func DeleteDatabase(w http.ResponseWriter, r *http.Request) {
	var database Database
	err := json.NewDecoder(r.Body).Decode(&database)
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

	err = db.Delete([]byte(database.Name), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Database " + database.Name + " successfully deleted")
}
