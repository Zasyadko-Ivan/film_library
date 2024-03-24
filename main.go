package main

import (
	"film_library/handler"
	"film_library/storage/postgres"
	"fmt"
	"log"
	"net/http"
)

const (
	host     = "db"
	port     = 5432
	user     = "root"
	password = "root"
	dbname   = "root"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := postgres.ConnectToDB(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("[INF] connection to the database is successful")

	err = db.Init()
	if err != nil {
		log.Fatal(err)
	}
	appHandler := handler.AppHandler{DB: db}

	http.HandleFunc("/actor/create", appHandler.AdminCheckMiddleware(appHandler.AddActor))
	http.HandleFunc("/actor/delete", appHandler.AdminCheckMiddleware(appHandler.DeleteActor))
	http.HandleFunc("/actor/change", appHandler.AdminCheckMiddleware(appHandler.ChangeActor))
	http.HandleFunc("/film/create", appHandler.AdminCheckMiddleware(appHandler.AddFilm))
	http.HandleFunc("/film/change", appHandler.AdminCheckMiddleware(appHandler.ChangeFilm))
	http.HandleFunc("/film/delete", appHandler.AdminCheckMiddleware(appHandler.DeleteFilm))
	http.HandleFunc("/film/add/actors", appHandler.AdminCheckMiddleware(appHandler.AddFilmActors))
	http.HandleFunc("/film/remove/actors", appHandler.AdminCheckMiddleware(appHandler.DeleteFilmActors))
	http.HandleFunc("/films", appHandler.UserCheckMiddleware(appHandler.GetAllFilms))
	http.HandleFunc("/films/name/film", appHandler.UserCheckMiddleware(appHandler.GetFilmsByNameFilm))
	http.HandleFunc("/films/name/actor", appHandler.UserCheckMiddleware(appHandler.GetFilmsByNameActor))
	http.HandleFunc("/actors", appHandler.UserCheckMiddleware(appHandler.GetAllActorOutFilms))

	http.ListenAndServe(":5000", nil)

}
