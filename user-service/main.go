package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	ID	int			`json:"id"`
	Name	string	`json:"name"`
	Email	string	`json:"email"`
}

var db *sql.DB

func main(){
	db.initDB()
	
	http.HandleFunc("/users",usersHandler)
	http.HandleFunc("/users/",userHandler)

	log.Println("User service is running on port 8081 .....")
	if err := http.ListenAndServe(":8081",nil); err != nil {
		log.Fatal("Server Failed :", err)
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case "GET":
		rows, err := db.Query("SELECT id,name,email FROM users")
		if err != nil{
			http.Error(w,"Database Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next(){
			var usr User
			if err:= rows.Scan(&usr.ID,&usr.Name,&usr.Email);err != nil{
				http.Error(w, "Error Scanning Data", http.StatusInternalServerError)
				return
			}
			users = append(users, usr)
		}
		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(users)
	
	case "POST":
		var usr User
		if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
			http.Error(w, "Invalid Request Payload", http.StatusInternalServerError)
			return
		}
		result, err := db.Exec("INSERT INTO users(name, email) VALUES(?, ?)",usr.Name, usr.Email)
		if err != nil {
			http.Error(w, "Failed to Insert User", http.StatusInternalServerError)
			return
		}
		id, _ = result.LastInsertId()
		usr.ID = int(id)
		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(usr)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request){
	idStr := r.URL.Path[len("/users/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	switch r.Method{
	case "GET":
		var usr User
		err := db.QueryRow("Select id, name, email FROM users WHERE id = ?", id).Scan(&usr.ID,&usr.Name,&usr.Email)
		if err != nil{
			http.Error(w, "User not Found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(usr)

	case "PUT":
		var usr User
		if err := json.NewDecoder(r.Body).Decode(&usr); err != nil{
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
			return
		}
		_,err := db.Exec("UPDATE Users SET name = ?, email = ? WHERE id = ?",usr.Name,usr.Name,id)
		if err != nil{
			http.Error(w, "Failed to Update User", http.StatusInternalServerError)
			return
		}
		usr.ID = id
		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(usr)

	case "DELETE":
		_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
		if err != nil{
			http.Error(w, "Failed to Delete User", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w,"Method Not Allowed", http.StatusMethodNotAllowed)
	}
}