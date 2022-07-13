package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"
)

type UpdateEmployeeInfo struct {
	Status          int            `json:"status,omitempty"`
	UpdatedEmployee store.Employee `json:"updatedEmployee,omitempty"`
}

//Pick up here where you left off.
func UpdateEmployee(w http.ResponseWriter, req *http.Request) {

	var updatedEmployee store.Employee

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&updatedEmployee)
	fmt.Println("JSON ERROR DECODING", err)

	//reqAdminID := req.Header["Admin-Id"][0]
	params := mux.Vars(req)
	employeeID := params["id"]

	//auditAction := fmt.Sprintf("User %v attemped to update wthe tags for working group %v to %v", reqAdminID, workingGroupID, updatedTag)

	if err != nil {
		//errorMsg := "Unknown field(s) included in request body or empty request body.  Please only editable employee information."
		//helpers.Panic(w, http.StatusBadRequest, errorMsg, err)
		fmt.Println("Unknown field(s) included in request body or empty request body.  Please only send editable employee information.")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "Unknown field(s) included in request body or empty request body.  Please only editable employee information."})
		return
	}
	//params := mux.Vars(r)
	//id := params["id"]

	validID, employeeErr := uuid.Parse(employeeID)

	//fmt.Println("YOUR VALID UUID", validID)

	// 400 bad request, invalid uuid
	if employeeErr != nil {
		fmt.Println("Status 400, not a valid uuid!")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "not a valid uuid!"})
		return
	}

	//fmt.Println("STORE EMPLOYEES", store.Employees)
	currentEmployeeInfo, ok := store.Employees[validID]
	// 404 not found, uuid does not exist
	if !ok {
		fmt.Println("Status 404, employee id not found.")
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "Employee ID not found!"})
		return
	}

	//else update the employee with the new fields....everything except the UUID can change!
	// type Employee struct {
	// 	Name       string    `json:"name,omitempty"`
	// 	Email      string    `json:"email,omitempty"`
	// 	EmployeeId uuid.UUID `json:"employeeId,omitempty"`
	// 	City       string    `json:"city,omitempty"`
	// 	Address    string    `json:"address,omitempty"`
	// 	Department string    `json:"department,omitempty"`
	// }

	//this blocks changing the employeeUUID as well.
	newEmployeeCopy := store.Employee{
		Name:       updatedEmployee.Name,
		Email:      updatedEmployee.Email,
		EmployeeId: currentEmployeeInfo.EmployeeId,
		City:       updatedEmployee.City,
		Address:    updatedEmployee.Address,
		Department: updatedEmployee.Department,
	}

	//print out old employee info!
	fmt.Println("BEFORE UPDATING INFO: ", currentEmployeeInfo)

	fmt.Println("AFTER UPDATING INFO: ", newEmployeeCopy)

	store.Employees[validID] = newEmployeeCopy

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	resp := UpdateEmployeeInfo{Status: 200, UpdatedEmployee: store.Employees[validID]}

	//helpers.MakeHistory(validWorkingGroupID, validAdminID, auditAction, http.StatusOK, "Success")
	json.NewEncoder(w).Encode(resp)

	fmt.Println("Successful update of employee info")
}
