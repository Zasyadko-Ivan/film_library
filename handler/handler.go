package handler

import (
	"encoding/json"
	"film_library/lib/e"
	"film_library/storage"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
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
		log.Printf("[ERR] [%d] %s", logNumber, MustPOST)
		http.Error(w, MustPOST, http.StatusMethodNotAllowed)
		return
	}

	var actor storage.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		log.Print(e.Wrap(logNumber, "failed to decode json", err))
		http.Error(w, ServerError, http.StatusInternalServerError)
		return
	}

	if actor.Name == "" || actor.Gender == "" || actor.Birthday == "" {
		log.Print(e.Wrap(logNumber, MissingRequiredFieldsActor, nil))
		http.Error(w, MissingRequiredFieldsActor, http.StatusBadRequest)
		return
	}

	_, err := time.Parse("2006-01-02", actor.Birthday)
	if err != nil {
		log.Print(e.Wrap(logNumber, InvalidBirthday, err))
		http.Error(w, InvalidBirthday, http.StatusBadRequest)
		return
	}
	log.Printf("[INF] [%d] start of the function execution AddActor", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution AddActor", logNumber)
	if err := ah.DB.AddActor(actor, logNumber); err == storage.ErrActorCreated {
		log.Print(e.Wrap(logNumber, "can't add actor to database", storage.ErrActorCreated))
		http.Error(w, ErrActorCreated, http.StatusConflict)
		return
	} else if err != nil {
		log.Print(err.Error())
		http.Error(w, ServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, ActorCreated)
}

func (ah *AppHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'DeleteActor' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'DeleteActor' handler has finished executing", logNumber)

	if r.Method != http.MethodDelete {
		log.Printf("[ERR] [%d] %s", logNumber, MustDELETE)
		http.Error(w, MustDELETE, http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	actorID := parts[len(parts)-1]

	log.Printf("[INF] [%d] start of the function execution DeleteActor actorID = %s", logNumber, actorID)
	defer log.Printf("[INF] [%d] end of the function execution DeleteActor", logNumber)
	if err := ah.DB.DeleteActor(actorID, logNumber); err != nil {
		log.Print(err.Error())
		http.Error(w, ServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, ActorDeleted)
}

func (ah *AppHandler) ChangeActor(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'ChangeActor' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'ChangeActor' handler has finished executing", logNumber)

	if r.Method != http.MethodPut {
		log.Printf("[ERR] [%d] %s", logNumber, MustPUT)
		http.Error(w, MustPUT, http.StatusMethodNotAllowed)
		return
	}

	var actor storage.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		log.Print(e.Wrap(logNumber, "failed to decode json", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if actor.Name == "" || actor.Gender == "" || actor.Birthday == "" {
		log.Print(e.Wrap(logNumber, MissingRequiredFieldsActor, nil))
		http.Error(w, MissingRequiredFieldsActor, http.StatusBadRequest)
		return
	}

	_, err := time.Parse("2006-01-02", actor.Birthday)
	if err != nil {
		log.Print(e.Wrap(logNumber, InvalidBirthday, err))
		http.Error(w, InvalidBirthday, http.StatusBadRequest)
		return
	}

	if actor.ReplaceBirthday != "" {
		_, err = time.Parse("2006-01-02", actor.ReplaceBirthday)
		if err != nil {
			log.Print(e.Wrap(logNumber, InvalidReleaseBirthday, err))
			http.Error(w, InvalidReleaseBirthday, http.StatusBadRequest)
			return
		}
	}

	log.Printf("[INF] [%d] start of the function execution ChangeActor", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution ChangeActor", logNumber)
	if err := ah.DB.ChangeActor(actor, logNumber); err == storage.ErrActorNotCreated {
		log.Print(e.Wrap(logNumber, "the actor is not in the database", storage.ErrActorNotCreated))
		http.Error(w, ErrActorNotCreated, http.StatusBadRequest)
		return
	} else if err == storage.ErrActorCreated {
		log.Print(e.Wrap(logNumber, "can't change actor to database", storage.ErrActorCreated))
		http.Error(w, ErrActorCreated, http.StatusConflict)
		return
	} else if err != nil {
		log.Print(err.Error())
		http.Error(w, ServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, ActorChange)
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

func (ah *AppHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'DeleteFilm' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'DeleteFilm' handler has finished executing", logNumber)

	if r.Method != http.MethodDelete {
		log.Printf("[ERR] [%d] The request method must be DELETE", logNumber)
		http.Error(w, "The request method must be DELETE", http.StatusMethodNotAllowed)
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

	log.Printf("[INF] [%d] start of the function execution DeleteFilm", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution DeleteFilm", logNumber)
	if err := ah.DB.DeleteFilm(film, logNumber); err == storage.ErrFilmNotCreated {
		log.Print(e.Wrap(logNumber, "the film is not in the database", storage.ErrFilmNotCreated))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Film successfully deleted")
}

func (ah *AppHandler) AddFilmActors(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'AddFilmActors' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'AddFilmActors' handler has finished executing", logNumber)

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

	log.Printf("[INF] [%d] start of the function execution AddFilmActors", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution AddFilmActors", logNumber)

	if err := ah.DB.AddFilmActors(film.Actors, film, logNumber); err == storage.ErrFilmNotCreated {
		log.Print(e.Wrap(logNumber, "the film is not in the database", storage.ErrFilmNotCreated))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "All the actors have been added to the film")
}

func (ah *AppHandler) DeleteFilmActors(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'DeleteFilmActors' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'DeleteFilmActors' handler has finished executing", logNumber)

	if r.Method != http.MethodDelete {
		log.Printf("[ERR] [%d] The request method must be DELETE", logNumber)
		http.Error(w, "The request method must be DELETE", http.StatusMethodNotAllowed)
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

	log.Printf("[INF] [%d] start of the function execution DeleteFilmActors", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution DeleteFilmActors", logNumber)

	if err := ah.DB.DeleteFilmActors(film.Actors, film, logNumber); err == storage.ErrFilmNotCreated {
		log.Print(e.Wrap(logNumber, "the film is not in the database", storage.ErrFilmNotCreated))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err == storage.ErrFilmNotActor {
		log.Print(e.Wrap(logNumber, "the film is not in the database", storage.ErrFilmNotActor))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "All the actors have been removed to the film")
}

func (ah *AppHandler) GetAllFilms(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'GetAllFilms' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'GetAllFilms' handler has finished executing", logNumber)

	if r.Method != http.MethodGet {
		log.Printf("[ERR] [%d] The request method must be GET", logNumber)
		http.Error(w, "The request method must be GET", http.StatusMethodNotAllowed)
		return
	}
	url := r.URL
	args := url.Query()

	sortByColoms := args.Get("sort_by_coloms")
	direction := args.Get("direction")

	log.Printf("[INF] [%d] start of the function execution GetAllFilms", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution GetAllFilms", logNumber)

	films, err := ah.DB.GetAllFilms(sortByColoms, direction, logNumber)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nameFilms := make(map[string][]string)
	nameFilms["name_films"] = films

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nameFilms)
}

func (ah *AppHandler) GetFilmsByNameFilm(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'GetFilmsByNameFilm' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'GetFilmsByNameFilm' handler has finished executing", logNumber)

	if r.Method != http.MethodGet {
		log.Printf("[ERR] [%d] The request method must be GET", logNumber)
		http.Error(w, "The request method must be GET", http.StatusMethodNotAllowed)
		return
	}
	url := r.URL
	args := url.Query()

	if !args.Has("name_film") {
		log.Print(e.Wrap(logNumber, "the 'name_film' argument is missing", nil))
		http.Error(w, "the 'name_film' argument is missing", http.StatusBadRequest)
		return
	}

	nameFilm := args.Get("name_film")

	log.Printf("[INF] [%d] start of the function execution GetFilmsByNameFilm", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution GetFilmsByNameFilm", logNumber)

	films, err := ah.DB.GetFilmsByNameFilm(nameFilm, logNumber)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nameFilms := make(map[string][]string)
	nameFilms["name_films"] = films

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nameFilms)
}

func (ah *AppHandler) GetFilmsByNameActor(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'GetFilmsByNameActor' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'GetFilmsByNameActor' handler has finished executing", logNumber)

	if r.Method != http.MethodGet {
		log.Printf("[ERR] [%d] The request method must be GET", logNumber)
		http.Error(w, "The request method must be GET", http.StatusMethodNotAllowed)
		return
	}
	url := r.URL
	args := url.Query()

	if !args.Has("name_actor") {
		log.Print(e.Wrap(logNumber, "the 'name_actor' argument is missing", nil))
		http.Error(w, "the 'name_actor' argument is missing", http.StatusBadRequest)
		return
	}

	nameActor := args.Get("name_actor")

	log.Printf("[INF] [%d] start of the function execution GetFilmsByNameАctor", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution GetFilmsByNameАctor", logNumber)

	films, err := ah.DB.GetFilmsByNameАctor(nameActor, logNumber)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nameFilms := make(map[string][]string)
	nameFilms["name_films"] = films

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nameFilms)
}

func (ah *AppHandler) GetAllActorOutFilms(w http.ResponseWriter, r *http.Request) {
	logNumber := genRandomNumber()

	log.Printf("[INF] [%d] the 'GetFilmsByNameActor' handler has started to run", logNumber)
	defer log.Printf("[INF] [%d] the 'GetFilmsByNameActor' handler has finished executing", logNumber)

	if r.Method != http.MethodGet {
		log.Printf("[ERR] [%d] The request method must be GET", logNumber)
		http.Error(w, "The request method must be GET", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("[INF] [%d] start of the function execution GetAllActorOutFilms", logNumber)
	defer log.Printf("[INF] [%d] end of the function execution GetAllActorOutFilms", logNumber)

	actorsOutFilms, err := ah.DB.GetAllActorOutFilms(logNumber)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actorsOutFilms)
}

func (ah *AppHandler) AdminCheckMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logNumber := genRandomNumber()
		log.Printf("[INF] [%d] the 'AdminCheckMiddleware' middleware has started to run", logNumber)
		defer log.Printf("[INF] [%d] the 'AdminCheckMiddleware' middleware has finished executing", logNumber)
		token := extractTokenFromHeader(r)

		if token == "" {
			log.Print(e.Wrap(logNumber, TokenNotFound, nil))
			http.Error(w, TokenNotFound, http.StatusBadRequest)
			return
		}

		right, err := ah.DB.CheckRightFromDB(token, logNumber)
		if err != nil {
			log.Print(e.Wrap(logNumber, "the rigth is not in the database", err))
			http.Error(w, TokenInvalid, http.StatusUnauthorized)
			return
		}

		if right != "admin" {
			log.Print(e.Wrap(logNumber, OnlyAdmin, nil))
			http.Error(w, OnlyAdmin, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (ah *AppHandler) UserCheckMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logNumber := genRandomNumber()
		token := extractTokenFromHeader(r)
		log.Printf("[INF] [%d] the 'UserCheckMiddleware' middleware has started to run", logNumber)
		defer log.Printf("[INF] [%d] the 'UserCheckMiddleware' middleware has finished executing", logNumber)

		if token == "" {
			log.Print(e.Wrap(logNumber, TokenNotFound, nil))
			http.Error(w, TokenNotFound, http.StatusBadRequest)
			return
		}

		right, err := ah.DB.CheckRightFromDB(token, logNumber)
		if err != nil || (right != "admin" && right != "user") {
			log.Print(e.Wrap(logNumber, "the rigth is not in the database", err))
			http.Error(w, TokenInvalid, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func extractTokenFromHeader(r *http.Request) string {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		return ""
	}

	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
