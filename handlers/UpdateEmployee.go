package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"
	"net/http"
)

func (s Server) UpdateEmployee(w http.ResponseWriter, req *http.Request) {

	var updatedEmployee store.Employee

	decoder := json.NewDecoder(req.Body)

	decoder.DisallowUnknownFields()

	if decoderErr := decoder.Decode(&updatedEmployee); decoderErr != nil {
		store.SendError(w, http.StatusBadRequest,
			"PATCH: Unknown field(s) included in request body or empty request body. Please only use editable employee information.",
			decoderErr)
		return
	}

	//pull ID from URL && check for validity
	params := mux.Vars(req)
	employeeID := params["id"]
	validID, invalidEmployeeErr := uuid.Parse(employeeID)

	// 400 bad request, invalid uuid
	if invalidEmployeeErr != nil {
		store.SendError(w, http.StatusBadRequest, "PATCH: not a valid uuid!", invalidEmployeeErr)
		return
	}

	output, errGetEmp4Update := s.Db.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"EmployeeId": &dynamodb.AttributeValue{
				S: aws.String(validID.String()),
			},
		},
		TableName: aws.String(TableName),
	})

	if errGetEmp4Update != nil {
		store.SendError(w, http.StatusInternalServerError, "PATCH: Unable to find employee!", errGetEmp4Update)
		return
	}

	currentEmployeeInfo := store.Employee{}
	unmarshallErr := dynamodbattribute.UnmarshalMap(output.Item, &currentEmployeeInfo)

	fmt.Println("Current employee info:", currentEmployeeInfo)

	// Checking to see if the employee info can be unmarshalled from DynamoDb AV
	//"PATCH: Problem unmarshalling your employee from Dynamo."
	if unmarshallErr != nil {
		store.SendError(w, http.StatusInternalServerError, "PATCH: Employee ID, not found-- unmarshall-able!", unmarshallErr)
		return
	}

	//Fill in a new employee based on the fields from the new struct, etc.
	newEmployeeCopy := store.Employee{
		Name:       updatedEmployee.Name,
		Email:      updatedEmployee.Email,
		EmployeeId: currentEmployeeInfo.EmployeeId,
		City:       updatedEmployee.City,
		Address:    updatedEmployee.Address,
		Department: updatedEmployee.Department,
	}

	//Old employee info!
	fmt.Println("BEFORE UPDATING INFO: ", currentEmployeeInfo)
	//Patched employee info
	fmt.Println("AFTER UPDATING INFO: ", newEmployeeCopy)

	//check all struct fields
	if newEmployeeCopy.Name == "" {
		newEmployeeCopy.Name = currentEmployeeInfo.Name
	}
	if newEmployeeCopy.Email == "" {
		newEmployeeCopy.Email = currentEmployeeInfo.Email
	}
	if newEmployeeCopy.EmployeeId == "" {
		newEmployeeCopy.EmployeeId = currentEmployeeInfo.EmployeeId
	}
	if newEmployeeCopy.City == "" {
		newEmployeeCopy.City = currentEmployeeInfo.City
	}
	if newEmployeeCopy.Address == "" {
		newEmployeeCopy.Address = currentEmployeeInfo.Address
	}
	if newEmployeeCopy.Department == "" {
		newEmployeeCopy.Department = currentEmployeeInfo.Department
	}

	//Store updated employee into dynamoDB
	_, updateItemErr := s.Db.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":city":         {S: aws.String(newEmployeeCopy.City)},
			":address":      {S: aws.String(newEmployeeCopy.Address)},
			":employeeName": {S: aws.String(newEmployeeCopy.Name)},
			":department":   {S: aws.String(newEmployeeCopy.Department)},
			":email":        {S: aws.String(newEmployeeCopy.Email)},
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
				employeeName = :employeeName,
				department = :department,
				email = :email`),
	})

	//FAILED TO UPDATE YOUR EMPLOYEE RECORD in Dynamo!
	if updateItemErr != nil {
		store.SendError(w, http.StatusInternalServerError, "PATCH: FAILED TO UPDATE YOUR EMPLOYEE RECORD!", updateItemErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := GetOneEmployee{Status: 200, Employee: newEmployeeCopy}
	json.NewEncoder(w).Encode(resp)
	fmt.Println("Update Employee Successfully Completed")
}
