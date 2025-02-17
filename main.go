package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Yandex-Practicum/final-project/handlers"
	"github.com/Yandex-Practicum/final-project/middleware"
	"github.com/Yandex-Practicum/final-project/storage"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Some error occured. Err: %s", err)
	}
	db, err := storage.CreateDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	storage := storage.NewTaskStorage(db)
	web_server_port := os.Getenv("TODO_PORT")
	mux := chi.NewRouter()
	mux.Handle("/*", http.FileServer(http.Dir("./web")))
	mux.Post("/api/signin", handlers.HangdleLogin)
	mux.Get("/api/nextdate", handlers.NextData)

	mux.Post("/api/task", middleware.Auth(handlers.HandleAddTask(storage)))
	mux.Get("/api/task", middleware.Auth(handlers.HandleGetTask(storage)))
	mux.Put("/api/task", middleware.Auth(handlers.HandleEditTask(storage)))
	mux.Delete("/api/task", middleware.Auth(handlers.HandleDeleteTask(storage)))

	mux.Post("/api/task/done", middleware.Auth(handlers.HandleTaskDone(storage)))

	mux.Get("/api/tasks", middleware.Auth(handlers.HandleGetTasks(storage)))

	err = http.ListenAndServe(":"+web_server_port, mux)
	if err != nil {
		panic(err)
	}
}
