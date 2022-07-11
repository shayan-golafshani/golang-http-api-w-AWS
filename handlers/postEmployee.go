package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Employee struct {
	Name       string    `json:"name,omitempty"`
	Email      string    `json:"email,omitempty"`
	EmployeeId uuid.UUID `json:"employeeId,omitempty"`
	City       string    `json:"city,omitempty"`
	Address    string    `json:"address,omitempty"`
	Department string    `json:"department,omitempty"`
}

type GetOneEmployee struct {
	Status   int      `json:"status,omitempty"`
	Employee Employee `json:"employee,omitempty"`
}

func PostEmployee(w http.ResponseWriter, r *http.Request) {

	var newEmployeeUUID = uuid.New()
	fmt.Println("Employee UUID:", newEmployeeUUID)

	//decode the employee from what's posted.
	var post Employee

	err := json.NewDecoder(r.Body).Decode(&post)

	fmt.Println("Error in creating employee", err)

	if err != nil {
		fmt.Println("You've got request body issues")
	}

	if post.Name == "" {
		fmt.Println("Please enter an Employee Name")
	}

	newEmployee := Employee{
		Name:       post.Name,
		Email:      (post.Name + "@twilio.com"),
		EmployeeId: newEmployeeUUID,
		City:       post.City,
		Address:    post.Address,
		Department: post.Department,
	}

	fmt.Println("New recruit ->", newEmployee)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	resp := GetOneEmployee{http.StatusCreated, newEmployee}

	json.NewEncoder(w).Encode(resp)
}
