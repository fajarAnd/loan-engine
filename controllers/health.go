package controller

import (
	db "github.com/fajar-andriansyah/loan-engine/infrastructure/database"
	"net/http"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {

	if err := db.GetConn().Ping(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
