package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	//"github.com/sendgrid/golang-http-api-w-AWS/store"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"
)

type GetOneEmployee struct {
	Status   int            `json:"status,omitempty"`
	Employee store.Employee `json:"employee,omitempty"`
}

func PostEmployee(w http.ResponseWriter, r *http.Request) {

	var newEmployeeUUID = uuid.New()
	fmt.Println("Employee UUID:", newEmployeeUUID)

	//decode the employee from what's posted.
	var post store.Employee

	err := json.NewDecoder(r.Body).Decode(&post)

	fmt.Println("Error in creating employee", err)

	if err != nil {
		fmt.Println("You've got request body issues")
	}

	if post.Name == "" {
		fmt.Println("Please enter an Employee Name")
	}

	newEmployee := store.Employee{
		Name:       post.Name,
		Email:      (post.Name + "@twilio.com"),
		EmployeeId: newEmployeeUUID,
		City:       post.City,
		Address:    post.Address,
		Department: post.Department,
	}

	fmt.Println("New recruit ->", newEmployee)

	//store the employee data
	store.Employees[newEmployee.EmployeeId] = newEmployee

	//probably going to need to store this to DynamoDB

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	resp := GetOneEmployee{http.StatusCreated, newEmployee}

	json.NewEncoder(w).Encode(resp)
}
