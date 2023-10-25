package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func login_Admin(w http.ResponseWriter, r *http.Request) {
	var admin_decode Admin
	admin := Admin{Name: "Admin", Password: "12345"}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&admin_decode); err != nil {
		http.Error(w, "JSON decoding error", http.StatusBadRequest)
		return
	}

	if admin_decode.Name == admin.Name && admin_decode.Password == admin.Password {
		json.NewEncoder(w).Encode(r.Body)
		CreateTokenHandler(w, r, admin_decode)
	} else {
		http.Error(w, "wrong name or password", http.StatusBadRequest)
		return
	}

}

func getUsers_admin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get connect
	db := connect()
	defer db.Close()

	if verfied, claims := VerifyTokenHandler(w, r); verfied == true && claims["username"].(string) == "Admin" {
		var user []User
		if err := db.Model(&user).Select(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Returning products
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "Invalid token", http.StatusBadRequest)
	}

}

func deleteUser_admin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Get connect
	db := connect()
	defer db.Close()

	var user User

	params := mux.Vars(r)
	productId := params["id"] // id = user's mail

	if verfied, claims := VerifyTokenHandler(w, r); verfied == true && claims["username"].(string) == "Admin" {
		if _, err := db.Model(&user).Where("email = ?", productId).Delete(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return

		} else {
			fmt.Println("Delete action succeeded")
		}

	} else {
		http.Error(w, "Invalid token", http.StatusBadRequest)
	}

}
