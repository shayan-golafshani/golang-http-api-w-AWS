package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Router http.Handler
	Db     store.SubsetDynamoDb
}

func (s *Server) SetupRouter() {
	s.Router = GetRouter(s)
}

const (
	TableName = "Employees"
)

// GetRouter builds the mux router used by a server
func GetRouter(server *Server) *mux.Router {
	fmt.Println("getting router on Lambda")
	//ADD EMPLOYEES TO THE STORE... will need to store these bad bois in DynamoDB for persistence..
	//store.AddEmployees()

	//New mux.Router
	r := mux.NewRouter()
	//USE logging middleware, w method durations, time-stamps, & define request type
	//r.Use(CustomMiddleware)

	//Create employee (C)
	r.HandleFunc("/v1/employee/add", func(w http.ResponseWriter, r *http.Request) {
		server.PostEmployee(w, r)
	}).Methods("POST")

	//Get Employee via employeeID (R)
	r.HandleFunc("/v1/employee/{id}", func(w http.ResponseWriter, r *http.Request) {
		server.GetEmployee(w, r)
	}).Methods("GET")

	//Update employee (U)
	r.HandleFunc("/v1/employee/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		server.UpdateEmployee(w, r)
	}).Methods("PATCH")

	//Delete employee (D)
	r.HandleFunc("/v1/employee/delete", func(w http.ResponseWriter, r *http.Request) {
		server.DeleteEmployee(w, r)
	}).Methods("DELETE")

	return r
}

// Router builds the mux router for the app and returns a Server
func Router(opts ...func(*Server)) (*Server, error) {
	server, err := NewServer(opts...)
	if err != nil {
		return nil, err
	}

	r := GetRouter(server)

	server.Router = r
	return server, nil
}

func NewServer(opts ...func(*Server)) (*Server, error) {
	server := &Server{}

	for _, opt := range opts {
		opt(server)
	}
	return server, nil
}

//Adds URI LOGGING, ADDS duration of how long the call took, adds Header application type JSON
func CustomMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log request
		log.Println(r.RequestURI)

		//Add header key/value to all requests
		w.Header().Add("Content-Type", "application/json")

		//start the time for the request...
		start := time.Now()

		//TODO check that this still works
		//wait for the request
		next.ServeHTTP(w, r)
		//print the duration
		fmt.Println("duration_ms", time.Since(start).Milliseconds())

		//pass on the rest of the request
		next.ServeHTTP(w, r)
	})
}
