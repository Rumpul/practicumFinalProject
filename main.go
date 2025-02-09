package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Yandex-Practicum/final-project/db"
	"github.com/Yandex-Practicum/final-project/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Some error occured. Err: %s", err)
	}
	db, err := db.CreateDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	web_server_port := os.Getenv("TODO_PORT")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	mux.HandleFunc("/api/nextdate", handlers.NextData)
	err = http.ListenAndServe(web_server_port, mux)
	if err != nil {
		panic(err)
	}
}
