package models

import (
	_ "modernc.org/sqlite"
)

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Login struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
