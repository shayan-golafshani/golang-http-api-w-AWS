package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"code.hq.twilio.com/hatch/twilio-working-groups-service/helpers"
	"github.com/google/uuid"
)

type postReqBody struct {
	EmployeeName string
}

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

func PostOneEmployee(w http.ResponseWriter, r *http.Request) {

	var newEmployeeUUID = uuid.New()
	fmt.Println("Employee UUID:", newEmployeeUUID)

	//reqAdminID := r.Header["Admin-Id"][0]
	var post postReqBody

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		helpers.Panic(w, http.StatusBadRequest, "Request Body Config Issues", err)
		return
	}

	if post.EmployeeName == "" {
		helpers.Panic(w, http.StatusBadRequest, "Please enter an Employee Name", err)
		return
	}

	// currentID, err := uuid.Parse(reqAdminID)
	// if err != nil {
	// 	helpers.Panic(w, http.StatusBadRequest, "Admin-Id in Headers not valid UUID", err)
	// }

	// if !helpers.UserExists(ds.UserMap, currentID) {
	// 	helpers.Panic(w, http.StatusNotFound, "Admin-Id in header not a registered user", err)
	// }

	newEmployee := Employee{
		Name:       post.EmployeeName,
		Email:      (post.EmployeeName + "@twilio.com"),
		EmployeeId: newEmployeeUUID,
		City:       "",
		Address:    "",
		Department: "",
	}

	fmt.Println(newEmployee)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	resp := GetOneEmployee{http.StatusCreated, newEmployee}

	json.NewEncoder(w).Encode(resp)
}
