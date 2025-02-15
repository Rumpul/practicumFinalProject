package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Yandex-Practicum/final-project/db"
	"github.com/Yandex-Practicum/final-project/handlers"
	"github.com/Yandex-Practicum/final-project/middleware"
	"github.com/go-chi/chi/v5"
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
	mux := chi.NewRouter()
	mux.Handle("/*", http.FileServer(http.Dir("./web")))
	mux.Get("/api/nextdate", handlers.NextData)

	mux.Post("/api/task", middleware.Auth(handlers.HandleAddTask(db)))
	mux.Get("/api/task", middleware.Auth(handlers.HandleGetTask(db)))
	mux.Put("/api/task", middleware.Auth(handlers.HandleEditTask(db)))
	mux.Delete("/api/task", middleware.Auth(handlers.HandleDeleteTask(db)))

	mux.Post("/api/task/done", middleware.Auth(handlers.HandleTaskDone(db)))

	mux.Get("/api/tasks", middleware.Auth(handlers.HandleGetTasks(db)))

	err = http.ListenAndServe(web_server_port, mux)
	if err != nil {
		panic(err)
	}
}
