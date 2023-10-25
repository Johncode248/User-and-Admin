package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db := connect()
	defer db.Close()

	r := mux.NewRouter()

	// User functions
	r.HandleFunc("/create_account", createAccount).Methods("POST")
	r.HandleFunc("/login/user", loginUser).Methods("POST")
	r.HandleFunc("/login/user", getInfoUser).Methods("GET")
	r.HandleFunc("/update/user", updateUser).Methods("PUT")
	// Admin functions
	r.HandleFunc("/admin/login", login_Admin).Methods("POST")
	r.HandleFunc("/admin/get_users", getUsers_admin).Methods("GET")
	r.HandleFunc("/admin/delete/{id}", deleteUser_admin).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8066", r))
}
