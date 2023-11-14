package main

import (
	"time"
)

type User struct {
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Date_birth time.Time    `json:"date_birth"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Updated_at time.Time `json:"updated_at"`
}

type Admin struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
