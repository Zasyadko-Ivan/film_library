package handler

import (
	"encoding/json"
	"film_library/storage"
	"net/http"
	"time"
)

type AppHandler struct {
	DB storage.Storage
}

func (ah *AppHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	var actor storage.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err := time.Parse("02-01-2006", actor.Birthday)
	if err != nil {
		http.Error(w, "Invalid birthday format. Please use DD-MM-YYYY", http.StatusBadRequest)
		return
	}

	if err := ah.DB.AddActor(actor); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(actor)
}