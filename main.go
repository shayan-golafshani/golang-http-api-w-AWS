package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/handlers"
)

func main() {

	myRouter := mux.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	fmt.Println("server starting on port: ", port)

	myRouter.HandleFunc("/v1/employee/add", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostOneEmployee(w, r)
	}).Methods("POST")

	err := http.ListenAndServe(port, myRouter)
	log.Fatal(err)
}
