package handler

import (
	"encoding/json"
	"film_library/lib/e"
	"film_library/storage"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type AppHandler struct {
	DB storage.Storage
}

func genRandomNumber() int {
	return rand.Intn(90000) + 10000
}

func (ah *AppHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'AddActor' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'AddActor' handler has finished executing", logNumber)
	var actor storage.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		log.Print(e.Wrap(logNumber, "failed to decode json", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if actor.Name == "" || actor.Gender == "" || actor.Birthday == "" {
		log.Print(e.Wrap(logNumber, "Required fields (name, gender, birthday) are missing", nil))
		http.Error(w, "Required fields (name, gender, birthday) are missing", http.StatusBadRequest)
		return
	}

	_, err := time.Parse("02-01-2006", actor.Birthday)
	if err != nil {
		log.Print(e.Wrap(logNumber, "Invalid birthday format. Please use DD-MM-YYYY", err))
		http.Error(w, "Invalid birthday format. Please use DD-MM-YYYY", http.StatusBadRequest)
		return
	}
	log.Printf("[INF] [%d] start of the function execution AddActor", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution AddActor", logNumber)
	if err := ah.DB.AddActor(actor, logNumber); err == storage.ErrActorCreated {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)
}
