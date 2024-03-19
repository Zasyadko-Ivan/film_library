package storage

import "errors"

type Actor struct {
	Name            string `json:"name_actor"`
	Gender          string `json:"gender"`
	Birthday        string `json:"birthday"`
	ReplaceName     string `json:"replace_name_actor"`
	ReplaceGender   string `json:"replace_gender"`
	ReplaceBirthday string `json:"replace_birthday"`
}

type Film struct {
	Actors             []Actor `json:"actors"`
	Name               string  `json:"name_film"`
	Description        string  `json:"description"`
	Rating             float32 `json:"rating"`
	ReleaseDate        string  `json:"release_date"`
	ReplaceName        string  `json:"replace_name_film"`
	ReplaceDescription string  `json:"replace_description"`
	ReplaceRating      float32 `json:"replace_rating"`
	ReplaceReleaseDate string  `json:"replace_release_date"`
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
	GetAllFilms(sortBy string) ([]string, error)
	GetFilmsByNameFilm(nameFilm string) ([]string, error)
	GetFilmsByNameАctor(nameActor string) ([]string, error)
}

var ErrActorCreated = errors.New("the actor has already been created")
var ErrActorNotCreated = errors.New("the actor is not created")
var ErrFilmCreated = errors.New("the film has already been created")
var ErrFilmNotCreated = errors.New("the film is not created")
