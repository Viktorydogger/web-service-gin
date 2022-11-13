package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "123123"
	DB_NAME     = "AvitoTest"
)

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	DB, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return DB
}

type Customer struct {
	CustomerID string `json:"custid"`
	Balance    string `json:"balance"`
}

type JsonResponse struct {
	Type    string     `json:"type"`
	Data    []Customer `json:"data"`
	Message string     `json:"message"`
}

func main() {
	// Init the mux router
	router := mux.NewRouter()
	// Route handles & endpoints
	// Get all movies
	router.HandleFunc("/customers/", GetCustomers).Methods("GET")
	// Create a movie
	router.HandleFunc("/customers/", CreateCustomer).Methods("POST")
	// Delete a specific movie by the movieID
	router.HandleFunc("/customers/{custid}", DeleteCustomer).Methods("DELETE")
	// Delete all movies
	router.HandleFunc("/customers/", DeleteCustomers).Methods("DELETE")
	// serve the app
	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8000", router))
}
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func GetCustomers(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting customers...")
	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM movies")
	// check errors
	checkErr(err)
	// var response []JsonResponse
	var customers []Customer

	// Foreach movie
	for rows.Next() {
		var id int
		var custId string
		var balancee string
		err = rows.Scan(&id, &custId, &balancee)
		// check errors
		checkErr(err)
		customers = append(customers, Customer{CustomerID: custId, Balance: balancee})
	}
	var response = JsonResponse{Type: "success", Data: customers}
	json.NewEncoder(w).Encode(response)
}
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	CustomerID := r.FormValue("custid")
	Balance := r.FormValue("balance")

	var response = JsonResponse{}

	if CustomerID == "" || Balance == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID or movieName parameter."}
	} else {
		db := setupDB()
		printMessage("Inserting customer into DB")
		fmt.Println("Inserting new customer with ID: " + CustomerID + " and balance: " + Balance)
		var lastInsertID int
		err := db.QueryRow("INSERT INTO movies(movieID, movieName) VALUES($1, $2) returning id;", CustomerID, Balance).Scan(&lastInsertID)
		// check errors
		checkErr(err)
		response = JsonResponse{Type: "success", Message: "The movie has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	CustomerID := params["custid"]
	var response = JsonResponse{}
	if CustomerID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
	} else {
		db := setupDB()
		printMessage("Deleting movie from DB")
		_, err := db.Exec("DELETE FROM movies where movieID = $1", CustomerID)
		// check errors
		checkErr(err)
		response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}
func DeleteCustomers(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Deleting all movies...")
	_, err := db.Exec("DELETE FROM movies")
	// check errors
	checkErr(err)
	printMessage("All movies have been deleted successfully!")
	var response = JsonResponse{Type: "success", Message: "All movies have been deleted successfully!"}
	json.NewEncoder(w).Encode(response)
}
