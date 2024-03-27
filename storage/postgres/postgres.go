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

var start int = 100000

func ConnectToDB(connStr string) (*Storage, error) {
	log.Print("[INF] getting started connecting to the table")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, e.Wrap(start, "can't open database", err)
	}
	if err := db.Ping(); err != nil {
		return nil, e.Wrap(start, "can't connect to database", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Init() error {
	log.Print("[INF] the beginning of the initialization of the table")

	q := `CREATE TABLE IF NOT EXISTS actors (id SERIAL PRIMARY KEY, name VARCHAR(50), gender VARCHAR(2), birthday DATE);
		CREATE TABLE IF NOT EXISTS films (id SERIAL PRIMARY KEY, name VARCHAR(150), description VARCHAR(1000), released DATE, rating DECIMAL(5,2), list_actors INTEGER[]);
		CREATE TABLE IF NOT EXISTS user_rights (id SERIAL PRIMARY KEY, user_right VARCHAR(20), token VARCHAR(30))`

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap(start, "can't create table", err)
	}

	return nil
}

func (s *Storage) AddActor(actor storage.Actor, logNumber int) error {
	exists, err := s.checkActor(actor.Name, actor.Gender, actor.Birthday, logNumber)
	if err != nil {
		return err
	}

	if exists {
		return storage.ErrActorCreated
	}

	q := `INSERT INTO actors (name, gender, birthday) VALUES($1, $2, $3);`
	if _, err := s.db.Exec(q, actor.Name, actor.Gender, actor.Birthday); err != nil {
		return e.Wrap(logNumber, "can't add actor to database", err)
	}
	return nil
}

func (s *Storage) ChangeActor(actor storage.Actor, logNumber int) error {
	exists, err := s.checkActor(actor.Name, actor.Gender, actor.Birthday, logNumber)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrActorNotCreated
	}

	if actor.ReplaceName == "" {
		actor.ReplaceName = actor.Name
	}
	if actor.ReplaceGender == "" {
		actor.ReplaceGender = actor.Gender
	}
	if actor.ReplaceBirthday == "" {
		actor.ReplaceBirthday = actor.Birthday
	}

	exists, err = s.checkActor(actor.ReplaceName, actor.ReplaceGender, actor.ReplaceBirthday, logNumber)
	if err != nil {
		return err
	}

	if exists {
		return storage.ErrActorCreated
	}

	q := `UPDATE actors SET name = $1, gender = $2, birthday = $3 WHERE name = $4 AND gender = $5 AND birthday = $6;`
	if _, err := s.db.Exec(q, actor.ReplaceName, actor.ReplaceGender, actor.ReplaceBirthday, actor.Name, actor.Gender, actor.Birthday); err != nil {
		return e.Wrap(logNumber, "can't update actors to database", err)
	}

	return nil
}

func (s *Storage) DeleteActor(actorID string, logNumber int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return e.Wrap(logNumber, "can't start transaction", err)
	}

	q := `UPDATE films SET list_actors = array_remove(list_actors, $1)`
	if _, err := tx.Exec(q, actorID); err != nil {
		tx.Rollback()
		return e.Wrap(logNumber, "can't update films table", err)
	}

	q = `DELETE FROM actors WHERE id = $1`
	if _, err := tx.Exec(q, actorID); err != nil {
		tx.Rollback()
		return e.Wrap(logNumber, "can't delete actor from actors table", err)
	}

	if err := tx.Commit(); err != nil {
		return e.Wrap(logNumber, "can't commit transaction", err)
	}

	return nil
}

func (s *Storage) checkActor(name string, gender string, birthday string, logNumber int) (bool, error) {
	q := `SELECT id FROM actors WHERE name = $1 AND gender = $2 AND birthday = $3;`
	var id int

	err := s.db.QueryRow(q, name, gender, birthday).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, e.Wrap(logNumber, "can't check actor to database", err)
	}

	return true, nil
}

func (s *Storage) AddFilm(film storage.Film, logNumber int) error {
	exists, err := s.checkFilm(film.Name, film.ReleaseDate, logNumber)
	if err != nil {
		return err
	}

	if exists {
		return storage.ErrFilmCreated
	}

	idActors := make([]int, 0, len(film.Actors))

	for i := 0; i < len(film.Actors); i++ {
		id, err := s.idActor(film.Actors[i].Name, film.Actors[i].Gender, film.Actors[i].Birthday, logNumber)
		if err != nil {
			return err
		}
		idActors = append(idActors, id)
	}

	q := `INSERT INTO films (name, description, released, rating, list_actors) VALUES($1, $2, $3, $4, $5);`

	if _, err := s.db.Exec(q, film.Name, film.Description, film.ReleaseDate, film.Rating, pq.Array(idActors)); err != nil {
		return e.Wrap(logNumber, "can't add film to database", err)
	}
	return nil
}

func (s *Storage) ChangeFilm(film storage.Film, logNumber int) error {
	exists, err := s.checkFilm(film.Name, film.ReleaseDate, logNumber)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrFilmNotCreated
	}

	exists, err = s.checkFilm(film.Name, film.ReplaceReleaseDate, logNumber)
	if err != nil {
		return err
	}

	if exists {
		return storage.ErrFilmCreated
	}

	if film.ReplaceReleaseDate == "" {
		film.ReplaceReleaseDate = film.ReleaseDate
	}
	if film.ReplaceDescription == "" {
		film.ReplaceDescription = film.Description
	}
	if film.ReplaceRating == "" {
		film.ReplaceRating = film.Rating
	}

	q := `UPDATE films SET description = $1, released = $2, rating = $3 WHERE name = $4 AND released = $5;`
	if _, err := s.db.Exec(q, film.ReplaceDescription, film.ReplaceReleaseDate, film.ReplaceRating, film.Name, film.ReleaseDate); err != nil {
		return e.Wrap(logNumber, "can't update films to database", err)
	}

	return nil
}

func (s *Storage) DeleteFilm(filmID string, logNumber int) error {
	q := `DELETE FROM films WHERE id = $1;`
	if _, err := s.db.Exec(q, filmID); err != nil {
		return e.Wrap(logNumber, "can't delete film to database", err)
	}
	return nil
}

func (s *Storage) checkFilm(name string, released string, logNumber int) (bool, error) {
	q := `SELECT id FROM films WHERE name = $1 AND released = $2;`
	var id int

	err := s.db.QueryRow(q, name, released).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, e.Wrap(logNumber, "can't check film to database", err)
	}

	return true, nil
}

func (s *Storage) idActor(name string, gender string, birthday string, logNumber int) (int, error) {
	q := `SELECT id FROM actors WHERE name = $1 AND gender = $2 AND birthday = $3;`
	var id int
	err := s.db.QueryRow(q, name, gender, birthday).Scan(&id)
	if err != nil {
		return id, e.Wrap(logNumber, "can't check actor to database", err)
	}

	return id, nil
}

func (s *Storage) AddFilmActors(actors []storage.Actor, film storage.Film, logNumber int) error {
	exists, err := s.checkFilm(film.Name, film.ReleaseDate, logNumber)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrFilmNotCreated
	}

	tx, err := s.db.Begin()

	if err != nil {
		return e.Wrap(logNumber, "can't start transaction", err)
	}

	for _, actor := range actors {
		id_actor, err := s.idActor(actor.Name, actor.Gender, actor.Birthday, logNumber)
		if err != nil {
			tx.Rollback()
			return err
		}
		var id_film int

		existQuery := `SELECT id FROM films WHERE $1 = ANY(list_actors) AND name = $2 AND released = $3;`
		err = tx.QueryRow(existQuery, id_actor, film.Name, film.ReleaseDate).Scan(&id_film)
		if err != sql.ErrNoRows {
			tx.Rollback()
			return nil
		}

		addQuery := `UPDATE films SET list_actors = array_append(list_actors, $1) WHERE name = $2 AND released = $3;`
		_, err = tx.Exec(addQuery, id_actor, film.Name, film.ReleaseDate)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return e.Wrap(logNumber, "can't commit transaction", err)
	}
	return nil
}

func (s *Storage) DeleteFilmActors(actors []storage.Actor, film storage.Film, logNumber int) error {
	exists, err := s.checkFilm(film.Name, film.ReleaseDate, logNumber)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrFilmNotCreated
	}

	tx, err := s.db.Begin()

	if err != nil {
		return e.Wrap(logNumber, "can't start transaction", err)
	}

	for _, actor := range actors {
		id_actor, err := s.idActor(actor.Name, actor.Gender, actor.Birthday, logNumber)
		if err != nil {
			tx.Rollback()
			return err
		}
		var id_film int

		existQuery := `SELECT id FROM films WHERE $1 = ANY(list_actors) AND name = $2 AND released = $3;`
		err = tx.QueryRow(existQuery, id_actor, film.Name, film.ReleaseDate).Scan(&id_film)
		if err == sql.ErrNoRows {
			tx.Rollback()
			return storage.ErrFilmNotActor
		}

		addQuery := `UPDATE films SET list_actors = array_remove(list_actors, $1) WHERE id = $2;`
		_, err = tx.Exec(addQuery, id_actor, id_film)
		if err != nil {
			tx.Rollback()
			return e.Wrap(logNumber, "the film is not in the database", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return e.Wrap(logNumber, "can't commit transaction", err)
	}

	return nil
}

func (s *Storage) GetAllFilms(sortByColoms, direction string, logNumber int) ([]string, error) {
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
		return nil, e.Wrap(logNumber, "error in the database when executing the 'GetAllFilms' query", err)
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, e.Wrap(logNumber, "error in the database when executing the 'GetAllFilms' query", err)
		}
		nameFilms = append(nameFilms, name)
	}

	return nameFilms, nil
}

func (s *Storage) GetFilmsByNameFilm(nameFilm string, logNumber int) ([]string, error) {
	var nameFilms []string

	q := `SELECT name FROM films WHERE name LIKE '%' || $1 || '%';`

	rows, err := s.db.Query(q, nameFilm)
	if err != nil {
		return nil, e.Wrap(logNumber, "error in the database when executing the 'GetFilmsByNameFilm' query", err)
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, e.Wrap(logNumber, "error in the database when executing the 'GetFilmsByNameFilm' query", err)
		}
		nameFilms = append(nameFilms, name)
	}

	return nameFilms, nil
}

func (s *Storage) GetFilmsByNameАctor(nameActor string, logNumber int) ([]string, error) {
	var nameFilms []string
	q := `SELECT DISTINCT f.name FROM films f JOIN actors a ON a.id = ANY(f.list_actors) WHERE a.name LIKE '%' || $1 || '%';`

	rows, err := s.db.Query(q, nameActor)
	if err != nil {
		return nil, e.Wrap(logNumber, "error in the database when executing the 'GetFilmsByNameАctor' query", err)
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, e.Wrap(logNumber, "error in the database when executing the 'GetFilmsByNameАctor' query", err)
		}
		nameFilms = append(nameFilms, name)
	}

	return nameFilms, nil
}

func (s *Storage) GetAllActorOutFilms(logNumber int) (map[string][]string, error) {
	actorFilms := make(map[string][]string)
	q := `
	SELECT a.name AS actor_name, f.name AS film_name
	FROM actors a
	LEFT JOIN films f ON a.id = ANY(f.list_actors)
	ORDER BY a.name;
	`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, e.Wrap(logNumber, "error in the database when executing the 'GetAllActorOutFilms' query", err)
	}

	for rows.Next() {
		var actorName, filmName sql.NullString
		if err := rows.Scan(&actorName, &filmName); err != nil {
			return nil, e.Wrap(logNumber, "error in the database when executing the 'GetAllActorOutFilms' query", err)
		}

		if filmName.Valid {
			actorFilms[actorName.String] = append(actorFilms[actorName.String], filmName.String)
		} else {
			actorFilms[actorName.String] = []string{""}
		}

	}

	return actorFilms, nil

}

func (s *Storage) CheckRightFromDB(token string, logNumber int) (string, error) {
	q := `SELECT user_right FROM user_rights WHERE token = $1;`
	var right string

	err := s.db.QueryRow(q, token).Scan(&right)

	if err != nil {
		return "", err
	}

	return right, nil

}
