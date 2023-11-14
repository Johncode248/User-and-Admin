package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)


func createAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get connect
	db := connect()
	defer db.Close()

	user_instance := &User{
		Updated_at: time.Now().Add(-24 * time.Hour),
	}
	// Decoding request
	_ = json.NewDecoder(r.Body).Decode(&user_instance)

	// converting password for database
	hashedPassword, err := hashPassword(user_instance.Password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user_instance.Password = hashedPassword
	
	//Inserting into database
	_, errInsert := db.Model(user_instance).Insert()
	if errInsert != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	// Returning json
	json.NewEncoder(w).Encode(user_instance)

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
	/*
	if user_decode.Name == user.Name && user_decode.Password == user.Password {
		json.NewEncoder(w).Encode(r.Body)
		CreateTokenHandler_user(w, r, user_decode)
	} else {
		http.Error(w, "Decoding error, data mismatch", http.StatusBadRequest)
		return
	}
        */
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user_decode.Password)); err != nil {
		http.Error(w, "Incorrect password or name", http.StatusBadRequest)
		return
	}
	CreateTokenHandler_user(w, r, user_decode)
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

		
		//json.NewEncoder(w).Encode(user.Name)
		//json.NewEncoder(w).Encode(user.Surname)
		//json.NewEncoder(w).Encode(user.Date_birth)   NOT JSON
		userInfo := map[string]interface{}{
			"name":       user.Name,
			"surname":    user.Surname,
			"date_birth": user.Date_birth,
		}
		// Returning JSON
		if err := json.NewEncoder(w).Encode(userInfo); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
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
		//json.NewEncoder(w).Encode(user) NOT JSON
		userInfo := map[string]interface{}{
			"name":       user.Name,
			"surname":    user.Surname,
			"date_birth": user.Date_birth,
			"email":      user.Email,
			"password":   user.Password,
		}
		// Returning JSON
		if err := json.NewEncoder(w).Encode(userInfo); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
func countUsers(db *sql.DB, table string) (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM " + table)
	err := row.Scan(&count)
	return count, err
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
