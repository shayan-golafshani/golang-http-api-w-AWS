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

func (s Server) UpdateEmployee(w http.ResponseWriter, req *http.Request) {

	var updatedEmployee store.Employee

	decoder := json.NewDecoder(req.Body)

	errGetEmp4Update := decoder.Decode(&updatedEmployee)
	if errGetEmp4Update != nil {
		fmt.Println("PATCH: Unknown field(s) included in request body or empty request body.  Please only send editable employee information.", errGetEmp4Update.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "PATCH: Unknown field(s) included in request body or empty request body. Please only use editable employee information."})
		return
	}

	params := mux.Vars(req)
	employeeID := params["id"]
	validID, invalidEmployeeErr := uuid.Parse(employeeID)

	// 400 bad request, invalid uuid
	if invalidEmployeeErr != nil {
		fmt.Println("PATCH: Status 400, not a valid uuid! : ", invalidEmployeeErr.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "PATCH: not a valid uuid!"})
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
		fmt.Println("PATCH: Error getting EmployeeID out of DynamoDB:", errGetEmp4Update.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "PATCH: Unable to find employee!"})
		return
	}

	currentEmployeeInfo := store.Employee{}
	unmarshallErr := dynamodbattribute.UnmarshalMap(output.Item, &currentEmployeeInfo)

	fmt.Println("PRINTING OUTPUT ->", output)

	// Checking to see if the employee info can be unmarshalled from DynamoDb AV
	if unmarshallErr != nil {
		fmt.Println("PATCH: Problem unmarshalling your employee from Dynamo.", unmarshallErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "PATCH: Employee ID, not found-- unmarshall-able!"})
		return
	}

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

	if updateItemErr != nil {
		fmt.Println("PATCH: FAILED TO UPDATE YOUR EMPLOYEE RECORD in Dynamo!", updateItemErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "PATCH: FAILED TO UPDATE YOUR EMPLOYEE RECORD!"})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	//reuse the same struct
	resp := GetEmployeeResponse{Status: 200, Data: newEmployeeCopy}
	json.NewEncoder(w).Encode(resp)
	fmt.Println("Successful update of employee info")
}
