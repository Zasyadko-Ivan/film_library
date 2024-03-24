package main

import (
	"film_library/config"
	"film_library/handler"
	"film_library/storage/postgres"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config := config.ReadConfig()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.DBname)
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
