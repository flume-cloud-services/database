package controllers

import (
	"encoding/json"
	"net/http"
	"time"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/dgrijalva/jwt-go"
)

type credentials struct {
	Name string `json:"name"`
	Secret string `json:"secret"`
}

type database struct {
	Name string
	Content string
}

// Login with the admin name and the token
func Login(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	admin := os.Getenv("FLUME_DATABASE_ADMIN")
	if len(admin) == 0 {
		admin = "admin"
	}
	secret := os.Getenv("FLUME_DATABASE_SECRET")
	if len(secret) == 0 {
		secret = "this_is_a_secret_token"
	}

	if creds.Name == admin && creds.Secret == secret {

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Username: creds.Name,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		json.NewEncoder(w).Encode("Success")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// GetData from database to JSON
func GetData(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

	tknStr := c.Value
	claims := &Claims{}

	secret := []byte(os.Getenv("FLUME_DATABASE_SECRET"))
	if len(secret) == 0 {
		secret = []byte("this_is_a_secret_token")
	}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := leveldb.OpenFile("level.db", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	var databases []database

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		newDb := database{Name: string(key), Content: string(value)}
		databases = append(databases, newDb)
	}
	iter.Release()
	err = iter.Error()

	json.NewEncoder(w).Encode(databases)
}