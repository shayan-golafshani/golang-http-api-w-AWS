package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "code.hq.twilio.com/hatch/twilio-working-groups-service/helpers"
	// "code.hq.twilio.com/hatch/twilio-working-groups-service/store"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	//
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"
)

type GetEmployeeResponse struct {
	Status int            `json:"status,omitempty"`
	Data   store.Employee `json:"data,omitempty"`
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GET EMPLOYEE CALLED")

	params := mux.Vars(r)
	id := params["id"]

	validID, err := uuid.Parse(id)

	//fmt.Println("YOUR VALID UUID", validID)

	// 400 bad request, invalid uuid
	if err != nil {
		fmt.Println("Status 400, not a valid uuid!")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "not a valid uuid!"})
		return
	}

	fmt.Println("STORE EMPLOYEES", store.Employees)
	_, ok := store.Employees[validID]
	// 404 not found, uuid does not exist
	if !ok {
		fmt.Println("Status 404, employee id not found.")
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "Employee ID not found!"})
		return
	}

	// 200 ok, uuid found
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	resp := GetEmployeeResponse{Status: 200, Data: store.Employees[validID]}

	json.NewEncoder(w).Encode(resp)
}
