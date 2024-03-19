package postgres

import (
	"database/sql"
	"film_library/lib/e"
	"film_library/storage"
	"log"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func ConnectToDB(connStr string) (*Storage, error) {
	log.Print("[INF] getting started connecting to the table")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}
	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Init() error {
	log.Print("[INF] the beginning of the initialization of the table")

	q := `CREATE TABLE IF NOT EXISTS actors (id SERIAL PRIMARY KEY, name VARCHAR(50), gender VARCHAR(2), birthday DATE);
		CREATE TABLE IF NOT EXISTS films (id SERIAL PRIMARY KEY, name VARCHAR(150), description VARCHAR(1000), released DATE, rating DECIMAL(5,2), list_actors INTEGER[]);`

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table", err)
	}

	return nil
}

func (s *Storage) AddActor(actor storage.Actor) error {
	log.Print("[INF] start of the function execution AddActor")
	exists, err := s.checkActor(actor.Name, actor.Gender, actor.Birthday)
	if err != nil {
		return err
	}

	if exists {
		return storage.ErrActorCreated
	}

	q := `INSERT INTO actors (name, gender, birthday) VALUES($1, $2, $3);`
	if _, err := s.db.Exec(q, actor.Name, actor.Gender, actor.Birthday); err != nil {
		return e.Wrap("can't add actor to database", err)
	}
	return nil
}

func (s *Storage) ChangeActor(actor storage.Actor) error {
	log.Print("[INF] start of the function execution ChangeActor")
	exists, err := s.checkActor(actor.Name, actor.Gender, actor.Birthday)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrActorNotCreated
	}
	log.Println(actor)
	if actor.ReplaceName == "" {
		actor.ReplaceName = actor.Name
	}
	if actor.ReplaceGender == "" {
		actor.ReplaceGender = actor.Gender
	}
	if actor.ReplaceBirthday == "" {
		actor.ReplaceBirthday = actor.Birthday
	}

	q := `UPDATE actors SET name = $1, gender = $2, birthday = $3 WHERE name = $4 AND gender = $5 AND birthday = $6;`
	if _, err := s.db.Exec(q, actor.ReplaceName, actor.ReplaceGender, actor.ReplaceBirthday, actor.Name, actor.Gender, actor.Birthday); err != nil {
		return e.Wrap("can't update actors to database", err)
	}
	return nil
}

func (s *Storage) DeleteActor(actor storage.Actor) error {
	log.Print("[INF] start of the function execution DeleteActor")
	exists, err := s.checkActor(actor.Name, actor.Gender, actor.Birthday)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrActorNotCreated
	}

	id, err := s.idActor(actor.Name, actor.Gender, actor.Birthday)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return e.Wrap("can't start transaction", err)
	}

	q := `UPDATE films SET list_actors = array_remove(list_actors, $1)`
	if _, err := tx.Exec(q, id); err != nil {
		tx.Rollback()
		return e.Wrap("can't update films table", err)
	}

	q = `DELETE FROM actors WHERE name = $1 AND gender = $2 AND birthday = $3`
	if _, err := tx.Exec(q, actor.Name, actor.Gender, actor.Birthday); err != nil {
		tx.Rollback()
		return e.Wrap("can't delete actor from actors table", err)
	}

	if err := tx.Commit(); err != nil {
		return e.Wrap("can't commit transaction", err)
	}

	return nil
}

func (s *Storage) checkActor(name string, gender string, birthday string) (bool, error) {
	q := `SELECT id FROM actors WHERE name = $1 AND gender = $2 AND birthday = $3;`
	var id int

	err := s.db.QueryRow(q, name, gender, birthday).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, e.Wrap("can't check actor to database", err)
	}

	return true, nil
}

func (s *Storage) AddFilm(film storage.Film) error {
	log.Print("[INF] start of the function execution AddFilm")
	exists, err := s.checkFilm(film.Name, film.ReleaseDate)
	if err != nil {
		return err
	}

	if exists {
		return storage.ErrFilmCreated
	}

	idActors := make([]int, 0, len(film.Actors))

	for i := 0; i < len(film.Actors); i++ {
		id, err := s.idActor(film.Actors[i].Name, film.Actors[i].Gender, film.Actors[i].Birthday)
		if err != nil {
			return err
		}
		idActors = append(idActors, id)
	}

	q := `INSERT INTO films (name, description, released, rating, list_actors) VALUES($1, $2, $3, $4, $5);`

	if _, err := s.db.Exec(q, film.Name, film.Description, film.ReleaseDate, film.Rating, pq.Array(idActors)); err != nil {
		return e.Wrap("can't add film to database", err)
	}
	return nil
}

func (s *Storage) ChangeFilm(film storage.Film) error {
	log.Print("[INF] start of the function execution ChangeFilm")
	exists, err := s.checkFilm(film.Name, film.ReleaseDate)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrFilmNotCreated
	}

	if film.ReplaceName == "" {
		film.ReplaceName = film.Name
	}
	if film.ReplaceReleaseDate == "" {
		film.ReplaceReleaseDate = film.ReleaseDate
	}
	if film.ReplaceDescription == "" {
		film.ReplaceDescription = film.Description
	}
	if film.ReplaceRating == 0 {
		film.ReplaceRating = film.Rating
	}

	q := `UPDATE films SET name = $1, description = $2, released = $3, rating = $4 WHERE name = $5 AND released = $6;`
	if _, err := s.db.Exec(q, film.ReplaceName, film.ReplaceDescription, film.ReplaceReleaseDate, film.ReplaceRating, film.Name, film.ReleaseDate); err != nil {
		return e.Wrap("can't update actors to database", err)
	}

	return nil
}

func (s *Storage) DeleteFilm(film storage.Film) error {
	log.Print("[INF] start of the function execution DeleteFilm")
	exists, err := s.checkFilm(film.Name, film.ReleaseDate)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrFilmNotCreated
	}

	q := `DELETE FROM films WHERE name = $1 AND released = $2;`
	if _, err := s.db.Exec(q, film.Name, film.ReleaseDate); err != nil {
		return e.Wrap("can't delete actor to database", err)
	}
	return nil
}

func (s *Storage) checkFilm(name string, released string) (bool, error) {
	q := `SELECT id FROM films WHERE name = $1 AND released = $2;`
	var id int

	err := s.db.QueryRow(q, name, released).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, e.Wrap("can't check film to database", err)
	}

	return true, nil
}

func (s *Storage) idActor(name string, gender string, birthday string) (int, error) {
	q := `SELECT id FROM actors WHERE name = $1 AND gender = $2 AND birthday = $3;`
	var id int
	err := s.db.QueryRow(q, name, gender, birthday).Scan(&id)
	if err != nil {
		return id, e.Wrap("can't check actor to database", err)
	}

	return id, nil
}

func (s *Storage) AddFilmActors(actor storage.Actor, film storage.Film) error {
	log.Print("[INF] start of the function execution AddFilmActors")
	exists, err := s.checkFilm(film.Name, film.ReleaseDate)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrFilmNotCreated
	}

	id_actor, err := s.idActor(actor.Name, actor.Gender, actor.Birthday)
	if err != nil {
		return err
	}
	var id_film int

	existQuery := `SELECT id FROM films WHERE $1 = ANY(list_actors) AND name = $2 AND released = $3;`
	err = s.db.QueryRow(existQuery, id_actor, film.Name, film.ReleaseDate).Scan(&id_film)
	if err != sql.ErrNoRows {
		return nil
	}

	addQuery := `UPDATE films SET list_actors = array_append(list_actors, $1) WHERE name = $2 AND released = $3;`
	_, err = s.db.Exec(addQuery, id_actor, film.Name, film.ReleaseDate)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteFilmActors(actor storage.Actor, film storage.Film) error {
	log.Print("[INF] start of the function execution DeleteFilmActors")
	exists, err := s.checkFilm(film.Name, film.ReleaseDate)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrFilmNotCreated
	}

	id_actor, err := s.idActor(actor.Name, actor.Gender, actor.Birthday)
	if err != nil {
		return err
	}
	var id_film int

	existQuery := `SELECT id FROM films WHERE $1 = ANY(list_actors) AND name = $2 AND released = $3;`
	err = s.db.QueryRow(existQuery, id_actor, film.Name, film.ReleaseDate).Scan(&id_film)
	if err == sql.ErrNoRows {
		return nil
	}

	addQuery := `UPDATE films SET list_actors = array_remove(list_actors, $1) WHERE id = $2;`
	_, err = s.db.Exec(addQuery, id_actor, id_film)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetAllFilms(sortByColoms, direction string) ([]string, error) {
	var nameFilms []string
	var q string
	if sortByColoms == "name" {
		q = `SELECT name FROM films ORDER BY name`
	} else if sortByColoms == "released" {
		q = `SELECT name FROM films ORDER BY released`
	} else {
		q = `SELECT name FROM films ORDER BY rating`
	}

	if direction == "ASC" {
		q += ` ASC;`
	} else {
		q += ` DESC;`
	}

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		nameFilms = append(nameFilms, name)
	}

	return nameFilms, nil
}
