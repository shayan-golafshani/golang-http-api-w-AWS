package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
func (s Server) UpdateEmployee(w http.ResponseWriter, req *http.Request) {

	var updatedEmployee store.Employee

	decoder := json.NewDecoder(req.Body)
	//decoder.DisallowUnknownFields()
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
	output, err := s.Db.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"EmployeeId": &dynamodb.AttributeValue{
				S: aws.String(validID.String()),
			},
		},
		TableName: aws.String(TableName),
	})
	if err != nil {
		//panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
		fmt.Println("UPDATE: Error getting EmployeeID out of DynamoDB:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "Unable to find employee!"})
		return
	}

	//currentEmployeeInfo, ok := store.Employees[validID]
	currentEmployeeInfo := store.Employee{}
	err2 := dynamodbattribute.UnmarshalMap(output.Item, &currentEmployeeInfo)

	fmt.Println("PRINTING OUTPUT.ITEM ->", output.Item)

	// Checking to see if the employee info can be unmarshalled from DynamoDb AV
	if err2 != nil {
		fmt.Println("UPDATE: Problem unmarshalling your employee from Dynamo.", err2.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "UPDATE: Employee ID, not found-- unmarshall-able!"})
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

	//store.Employees[validID] = newEmployeeCopy

	//Store updated employee into dynamoDB
	_, updateItemErr := s.Db.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":city":       {S: aws.String(newEmployeeCopy.City)},
			":address":    {S: aws.String(newEmployeeCopy.Address)},
			":name":       {S: aws.String(newEmployeeCopy.Name)},
			":department": {S: aws.String(newEmployeeCopy.Department)},
			":email":      {S: aws.String(newEmployeeCopy.Email)},
		},
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"EmployeeId": {
				S: aws.String(currentEmployeeInfo.EmployeeId),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String(
			`set city = :city,
				address = :address,
				name = :name,
				department = :department,
				email = :email`),
	})

	if updateItemErr != nil {
		fmt.Println("FAILED TO UPDATE YOUR EMPLOYEE RECORD!", updateItemErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "FAILED TO UPDATE YOUR EMPLOYEE RECORD!"})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := UpdateEmployeeInfo{Status: 200, UpdatedEmployee: newEmployeeCopy}
	json.NewEncoder(w).Encode(resp)
	fmt.Println("Successful update of employee info")
}
