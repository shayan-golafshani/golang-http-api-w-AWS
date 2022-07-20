package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

func (s Server) GetEmployee(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GET EMPLOYEE CALLED")

	params := mux.Vars(r)
	id := params["id"]

	validID, err := uuid.Parse(id)

	//fmt.Println("YOUR VALID UUID", validID)

	// 400 bad request, invalid uuid
	if err != nil {
		fmt.Println("Status 400, not a valid uuid!", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "not a valid uuid!"})
		return
	}

	//Get your employee for the get
	output, getEmployeeErr := s.Db.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"EmployeeId": &dynamodb.AttributeValue{
				S: aws.String(validID.String()),
			},
		},
		TableName: aws.String(TableName),
	})

	if getEmployeeErr != nil {
		fmt.Println("GET: Error getting EmployeeID out of DynamoDB:", getEmployeeErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "Unable to find employee! Try again"})
		return
	}

	employeeInfo := store.Employee{}
	unmarshallErr := dynamodbattribute.UnmarshalMap(output.Item, &employeeInfo)
	fmt.Println("STORE EMPLOYEES", output.Item)

	if unmarshallErr != nil {
		fmt.Println("GET: Problem unmarshalling your employee from Dynamo.", unmarshallErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "GET: Employee ID, not found-- unmarshall-able!"})
		return
	}

	//_, ok := store.Employees[validID]
	// 404 not found, uuid does not exist
	//if !ok {
	//	fmt.Println("Status 404, employee id not found.")
	//	w.WriteHeader(http.StatusNotFound)
	//	w.Header().Set("Content-Type", "application/json")
	//	json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "GET: Employee ID not found!"})
	//	return
	//}

	// 200 ok, uuid found
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	resp := GetEmployeeResponse{Status: 200, Data: employeeInfo}

	json.NewEncoder(w).Encode(resp)
}
