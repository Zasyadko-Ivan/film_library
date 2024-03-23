package handler

import (
	"encoding/json"
	"film_library/lib/e"
	"film_library/storage"
	"fmt"
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

	if r.Method != http.MethodPost {
		log.Printf("[ERR] [%d] The request method must be POST", logNumber)
		http.Error(w, "The request method must be POST", http.StatusMethodNotAllowed)
		return
	}

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

	_, err := time.Parse("2006-01-02", actor.Birthday)
	if err != nil {
		log.Print(e.Wrap(logNumber, "Invalid birthday format. Please use YYYY-MM-DD", err))
		http.Error(w, "Invalid birthday format. Please use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	log.Printf("[INF] [%d] start of the function execution AddActor", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution AddActor", logNumber)
	if err := ah.DB.AddActor(actor, logNumber); err == storage.ErrActorCreated {
		log.Print(e.Wrap(logNumber, "can't add actor to database", storage.ErrActorCreated))
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Actor successfully create")
}

func (ah *AppHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'DeleteActor' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'DeleteActor' handler has finished executing", logNumber)

	if r.Method != http.MethodDelete {
		log.Printf("[ERR] [%d] The request method must be DELETE", logNumber)
		http.Error(w, "The request method must be DELETE", http.StatusMethodNotAllowed)
		return
	}

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

	_, err := time.Parse("2006-01-02", actor.Birthday)
	if err != nil {
		log.Print(e.Wrap(logNumber, "Invalid birthday format. Please use YYYY-MM-DD", err))
		http.Error(w, "Invalid birthday format. Please use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	log.Printf("[INF] [%d] start of the function execution DeleteActor", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution DeleteActor", logNumber)
	if err := ah.DB.DeleteActor(actor, logNumber); err == storage.ErrActorNotCreated {
		log.Print(e.Wrap(logNumber, "the actor is not in the database", storage.ErrActorNotCreated))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Actor successfully deleted")
}

func (ah *AppHandler) ChangeActor(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'ChangeActor' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'ChangeActor' handler has finished executing", logNumber)

	if r.Method != http.MethodPut {
		log.Printf("[ERR] [%d] The request method must be PUT", logNumber)
		http.Error(w, "The request method must be PUT", http.StatusMethodNotAllowed)
		return
	}

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

	_, err := time.Parse("2006-01-02", actor.Birthday)
	if err != nil {
		log.Print(e.Wrap(logNumber, "Invalid birthday format. Please use YYYY-MM-DD", err))
		http.Error(w, "Invalid birthday format. Please use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	if actor.ReplaceBirthday != "" {
		_, err = time.Parse("2006-01-02", actor.ReplaceBirthday)
		if err != nil {
			log.Print(e.Wrap(logNumber, "Invalid birthday format. Please use YYYY-MM-DD", err))
			http.Error(w, "Invalid birthday format. Please use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	log.Printf("[INF] [%d] start of the function execution ChangeActor", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution ChangeActor", logNumber)
	if err := ah.DB.ChangeActor(actor, logNumber); err == storage.ErrActorNotCreated {
		log.Print(e.Wrap(logNumber, "the actor is not in the database", storage.ErrActorNotCreated))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err == storage.ErrActorCreated {
		log.Print(e.Wrap(logNumber, "can't change actor to database", storage.ErrActorCreated))
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Actor successfully change")
}

func (ah *AppHandler) AddFilm(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'AddFilm' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'AddFilm' handler has finished executing", logNumber)

	if r.Method != http.MethodPost {
		log.Printf("[ERR] [%d] The request method must be POST", logNumber)
		http.Error(w, "The request method must be POST", http.StatusMethodNotAllowed)
		return
	}

	var film storage.Film
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		log.Print(e.Wrap(logNumber, "failed to decode json", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if film.Name == "" || film.ReleaseDate == "" {
		log.Print(e.Wrap(logNumber, "Required fields (name, release_date) are missing", nil))
		http.Error(w, "Required fields (name, release_date) are missing", http.StatusBadRequest)
		return
	}

	_, err := time.Parse("2006-01-02", film.ReleaseDate)
	if err != nil {
		log.Print(e.Wrap(logNumber, "Invalid release date format. Please use YYYY-MM-DD", err))
		http.Error(w, "Invalid date format format. Please use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	log.Printf("[INF] [%d] start of the function execution AddFilm", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution AddFilm", logNumber)
	if err := ah.DB.AddFilm(film, logNumber); err == storage.ErrFilmCreated {
		log.Print(e.Wrap(logNumber, "can't add film to database", storage.ErrFilmCreated))
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Film successfully create")
}

func (ah *AppHandler) ChangeFilm(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'ChangeFilm' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'ChangeFilm' handler has finished executing", logNumber)

	if r.Method != http.MethodPut {
		log.Printf("[ERR] [%d] The request method must be PUT", logNumber)
		http.Error(w, "The request method must be PUT", http.StatusMethodNotAllowed)
		return
	}

	var film storage.Film
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		log.Print(e.Wrap(logNumber, "failed to decode json", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if film.Name == "" || film.ReleaseDate == "" {
		log.Print(e.Wrap(logNumber, "Required fields (name, release_date) are missing", nil))
		http.Error(w, "Required fields (name, release_date) are missing", http.StatusBadRequest)
		return
	}

	_, err := time.Parse("2006-01-02", film.ReleaseDate)
	if err != nil {
		log.Print(e.Wrap(logNumber, "Invalid birthday format. Please use YYYY-MM-DD", err))
		http.Error(w, "Invalid birthday format. Please use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	if film.ReplaceReleaseDate != "" {
		_, err = time.Parse("2006-01-02", film.ReplaceReleaseDate)
		if err != nil {
			log.Print(e.Wrap(logNumber, "Invalid birthday format. Please use YYYY-MM-DD", err))
			http.Error(w, "Invalid birthday format. Please use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	log.Printf("[INF] [%d] start of the function execution ChangeFilm", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution ChangeFilm", logNumber)
	if err := ah.DB.ChangeFilm(film, logNumber); err == storage.ErrFilmNotCreated {
		log.Print(e.Wrap(logNumber, "the film is not in the database", storage.ErrFilmNotCreated))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err == storage.ErrFilmCreated {
		log.Print(e.Wrap(logNumber, "can't change film to database", storage.ErrFilmCreated))
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Film successfully change")
}
