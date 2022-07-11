package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"
)

type deletionReq struct {
	EmployeeId uuid.UUID `json:"employeeId,omitempty"`
}

type deletionResponse struct {
	Status    int
	Employees map[uuid.UUID]store.Employee
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {

	var post deletionReq

	err := json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		//helpers.Panic(w, 400, "Request Body Config Issues", err)
		fmt.Println("Deletion Request Body Config Issues")
		return
	}

	if post.EmployeeId.String() == "" {
		fmt.Println("Please Submit EmployeeId for deletion")

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "Please Submit EmployeeId for deletion"})
		return
	}

	//PARSE BODY
	validID, err := uuid.Parse(post.EmployeeId.String())

	//CHECK FOR A VALID UUID
	_, ok := store.Employees[validID]

	if !ok {
		fmt.Println("Status 404, Employee ID not found.")

		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 404, Msg: "Employee Id not valid, can't be deleted"})
		return
	}

	if ok {
		delete(store.Employees, validID)
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	resp := deletionResponse{204, store.Employees}

	json.NewEncoder(w).Encode(resp)

	fmt.Println("Removed employee! さよなら！")
	fmt.Println("/n ------------------------------------- /n")
	fmt.Println(store.Employees)
}
