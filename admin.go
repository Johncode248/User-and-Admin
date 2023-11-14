package main

import (
	"encoding/json"
	"fmt"
	"strconv"

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
		//json.NewEncoder(w).Encode(r.Body)  Empty JSON
		CreateTokenHandler(w, r, admin_decode)
	} else {
		http.Error(w, "wrong name or password", http.StatusBadRequest)
		return
	}

}

type Page struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Users    []User `json:"users"`
}

func getUsers_admin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get connect
	db := connect()
	defer db.Close()

	if verfied, claims := VerifyTokenHandler(w, r); verfied == true && claims["username"].(string) == "Admin" {
		// Set default values
		page := 1
		pageSize := 10

		// Check if user provided values for page
		if r.URL.Query().Get("page") != "" {
			if val, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
				page = val
			}
		}

		// Fetch data from the database
		var users []User
		if err := db.Model(&users).Limit(pageSize).Offset((page - 1) * pageSize).Select(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Return the page data
		json.NewEncoder(w).Encode(Page{Page: page, PageSize: pageSize, Users: users})
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
