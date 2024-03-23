package storage

import "errors"

type Actor struct {
	Name            string `json:"name_actor"`
	Gender          string `json:"gender"`
	Birthday        string `json:"birthday"`
	ReplaceName     string `json:"replace_name_actor,omitempty"`
	ReplaceGender   string `json:"replace_gender,omitempty"`
	ReplaceBirthday string `json:"replace_birthday,omitempty"`
}

type Film struct {
	Actors             []Actor `json:"actors"`
	Name               string  `json:"name_film"`
	Description        string  `json:"description"`
	Rating             string  `json:"rating"`
	ReleaseDate        string  `json:"release_date"`
	ReplaceDescription string  `json:"replace_description,omitempty"`
	ReplaceRating      string  `json:"replace_rating,omitempty"`
	ReplaceReleaseDate string  `json:"replace_release_date,omitempty"`
}

type Storage interface {
	AddActor(actor Actor, logNumber int) error
	ChangeActor(actor Actor, logNumber int) error
	DeleteActor(actor Actor, logNumber int) error
	AddFilmActors(actor []Actor, film Film, logNumber int) error
	DeleteFilmActors(actor Actor, film Film, logNumber int) error
	AddFilm(film Film, logNumber int) error
	ChangeFilm(film Film, logNumber int) error
	DeleteFilm(film Film, logNumber int) error
	GetAllFilms(sortByColoms, direction string, logNumber int) ([]string, error)
	GetFilmsByNameFilm(nameFilm string, logNumber int) ([]string, error)
	GetFilmsByNameАctor(nameActor string, logNumber int) ([]string, error)
	GetAllActorOutFilms(logNumber int) (map[string][]string, error)
}

var ErrActorCreated = errors.New("the actor has already been created")
var ErrActorNotCreated = errors.New("the actor is not created")
var ErrFilmCreated = errors.New("the film has already been created")
var ErrFilmNotCreated = errors.New("the film is not created")
