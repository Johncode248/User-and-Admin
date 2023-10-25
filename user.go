package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func createAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get connect
	db := connect()
	defer db.Close()

	user_instance := &User{

		Updated_at: time.Now().Add(-24 * time.Hour),
	}

	//10 users maximum

	if count, err := db.Model(&User{}).Count(); err != nil {
		log.Fatal(err)
		fmt.Println(count)
	} else {
		if count >= 10 {
			http.Error(w, "Maximum number of accounts reached.", http.StatusBadRequest)

		} else {
			// Decoding request
			_ = json.NewDecoder(r.Body).Decode(&user_instance)

			//Inserting into database
			_, err := db.Model(user_instance).Insert()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}

			// Returning json
			json.NewEncoder(w).Encode(user_instance)
		}

	}

}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var user_decode User
	db := connect()
	defer db.Close()

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user_decode); err != nil {
		http.Error(w, "JSON decoding error", http.StatusBadRequest)
		return
	}

	var user User

	err := db.Model(&user).Where("name = ?", user_decode.Name).Select()
	if err != nil {
		http.Error(w, "Incorrect password or name", http.StatusBadRequest)
	}
	if user_decode.Name == user.Name && user_decode.Password == user.Password {
		json.NewEncoder(w).Encode(r.Body)
		CreateTokenHandler_user(w, r, user_decode)
	} else {
		http.Error(w, "Decoding error, data mismatch", http.StatusBadRequest)
		return
	}
}

func getInfoUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Get connect
	db := connect()
	defer db.Close()

	var user User

	if verfied, claims := VerifyTokenHandler(w, r); verfied == true {
		if err := db.Model(&user).Where("name = ?", claims["username"].(string)).Select(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Returning json
		json.NewEncoder(w).Encode(user.Name)
		json.NewEncoder(w).Encode(user.Surname)
		json.NewEncoder(w).Encode(user.Date_birth)
	} else {
		http.Error(w, "invalid token", http.StatusBadRequest)
	}

}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get connection to the database
	db := connect()
	defer db.Close()

	var user User

	// Decode request body into the user_info struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if verfied, claims := VerifyTokenHandler(w, r); verfied == true {
		// checking last update, edit data possible once every 24 hours
		if err := db.Model(&user).Where("name = ?", claims["username"].(string)).Select(user.Updated_at); user.Updated_at.After(time.Now().Add(-24 * time.Hour)) {
			// updating data
			_, err = db.Model(&user).Where("name = ?", claims["username"].(string)).Set("name = ?, surname = ?, date_birth = ?,updated_at = ?", user.Name, user.Surname, user.Date_birth, time.Now()).Update()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Too early update", http.StatusBadRequest)
			return
		}

		// Returning user information
		json.NewEncoder(w).Encode(user)
	}
}
func countUsers(db *sql.DB, table string) (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM " + table)
	err := row.Scan(&count)
	return count, err
}
