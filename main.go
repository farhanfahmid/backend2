package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Define the User struct to represent the User entity
type User struct {
	ID             int64  `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Country        string `json:"country"`
	ProfilePicture string `json:"profile_picture"`
}

var db *sql.DB

func init() {
	var err error
	// Set up a global database connection
	db, err = sql.Open("mysql", "farhanfahmid:rootpassword2@tcp(localhost:3306)/rankapi")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
}

// Handler function for retrieving the list of users
func getUsers(w http.ResponseWriter, r *http.Request) {
	// Perform the database query to get the list of users
	query := "SELECT * FROM users"
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the users
	var users []User

	// Iterate over the rows and populate the users slice
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Country, &user.ProfilePicture); err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Respond with the list of users
	json.NewEncoder(w).Encode(users)
}

func main() {
	// Set up the router
	router := mux.NewRouter()

	// Enable CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Define the API endpoints for creating, updating, and deleting users
	router.HandleFunc("/users", createUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", updateUser).Methods(http.MethodPatch)
	router.HandleFunc("/users/{id}", deleteUser).Methods(http.MethodDelete)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the User Management API"))
	})

	// Run the server
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(router)))
}

// Handler function for creating a user
func createUser(w http.ResponseWriter, r *http.Request) {
	// Get the user data from the request body
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform the database insert operation to create a new user
	insertQuery := "INSERT INTO users (first_name, last_name, country, profile_picture) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(insertQuery, user.FirstName, user.LastName, user.Country, user.ProfilePicture)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Get the ID of the created user
	user.ID, _ = result.LastInsertId()

	// Respond with the created user
	json.NewEncoder(w).Encode(user)
}

// Handler function for updating a user
func updateUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL parameter
	vars := mux.Vars(r)
	id := vars["id"]

	// Get the user data from the request body
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform the database update operation to update the user
	updateQuery := "UPDATE users SET first_name = ?, last_name = ?, country = ?, profile_picture = ? WHERE id = ?"
	_, err := db.Exec(updateQuery, user.FirstName, user.LastName, user.Country, user.ProfilePicture, id)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Respond with the updated user
	json.NewEncoder(w).Encode(user)
}

// Handler function for deleting a user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL parameter
	vars := mux.Vars(r)
	id := vars["id"]

	// Perform the database delete operation to delete the user
	deleteQuery := "DELETE FROM users WHERE id = ?"
	_, err := db.Exec(deleteQuery, id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
}
