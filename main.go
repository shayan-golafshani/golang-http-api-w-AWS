package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/handlers"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"
)

func main() {
	store.AddEmployees()

	myRouter := mux.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	fmt.Println("server starting on port: ", port)

	//Get employee details using EmployeeID
	myRouter.HandleFunc("/v1/employee/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetEmployee(w, r)
	}).Methods("GET")

	myRouter.HandleFunc("/v1/employee/add", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostEmployee(w, r)
	}).Methods("POST")

	err := http.ListenAndServe(port, myRouter)
	log.Fatal(err)
}
