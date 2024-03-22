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

	http.HandleFunc("/actors/create", appHandler.AddActor)
	http.HandleFunc("/actors/delete", appHandler.DeleteActor)
	http.HandleFunc("/actors/change", appHandler.ChangeActor)

	http.ListenAndServe(":5000", nil)

}
