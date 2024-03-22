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
	Rating             float32 `json:"rating"`
	ReleaseDate        string  `json:"release_date"`
	ReplaceName        string  `json:"replace_name_film,omitempty"`
	ReplaceDescription string  `json:"replace_description,omitempty"`
	ReplaceRating      float32 `json:"replace_rating,omitempty"`
	ReplaceReleaseDate string  `json:"replace_release_date,omitempty"`
}

type Storage interface {
	AddActor(actor Actor) error
	ChangeActor(actor Actor) error
	DeleteActor(actor Actor) error
	AddFilmActors(actor Actor, film Film) error
	DeleteFilmActors(actor Actor, film Film) error
	AddFilm(film Film) error
	ChangeFilm(film Film) error
	DeleteFilm(film Film) error
	GetAllFilms(sortByColoms, direction string) ([]string, error)
	GetFilmsByNameFilm(nameFilm string) ([]string, error)
	GetFilmsByNameАctor(nameActor string) ([]string, error)
	GetAllActorOutFilms() (map[string][]string, error)
}

var ErrActorCreated = errors.New("the actor has already been created")
var ErrActorNotCreated = errors.New("the actor is not created")
var ErrFilmCreated = errors.New("the film has already been created")
var ErrFilmNotCreated = errors.New("the film is not created")
