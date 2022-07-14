package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/handlers"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"

	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/sendgrid/mc-contacts-custom-fields/internal/handlers"
	"github.com/sendgrid/mcauto/apigw"
	"github.com/sendgrid/mclogger/lib/logger"
)

type Server struct {
	router http.Handler
}

// Router builds the mux router for the app and returns a Server
func Router(opts ...func(*Server)) (*Server, error) {
	server, err := NewServer(opts...)
	if err != nil {
		return nil, err
	}

	r := GetRouter(server)

	server.router = r
	return server, nil
}

func NewServer(opts ...func(*Server)) (*Server, error) {
	server := &Server{}

	for _, opt := range opts {
		opt(server)
	}
	return server, nil
}

// GetRouter builds the mux router used by a server
func GetRouter(server *Server) *mux.Router {
	fmt.Println("getting routert on Lambda")
	//ADD EMPLOYEES TO THE STORE... will need to store these bad bois in DynamoDB for persistence..
	store.AddEmployees()

	//New mux.Router
	r := mux.NewRouter()

	//Create employee (C)
	r.HandleFunc("/v1/employee/add", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostEmployee(w, r)
	}).Methods("POST")

	//Get Employee via employeeID (R)
	r.HandleFunc("/v1/employee/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetEmployee(w, r)
	}).Methods("GET")

	//Update employee (U)
	r.HandleFunc("/v1/employee/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateEmployee(w, r)
	}).Methods("PATCH")

	//Delete employee (D)
	r.HandleFunc("/v1/employee/delete", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteEmployee(w, r)
	}).Methods("DELETE")

	return r
}

func main() {
	handler, err := Router()
	if err != nil {
		logger.NewEntry().Fatal("Unable to create http router for lambda:httpServer", err)
	}
	awslambda.Start(apigw.Handle(handler))
}
