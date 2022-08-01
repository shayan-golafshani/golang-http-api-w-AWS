package handlers

import (
	"encoding/json"
	"errors"
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
	fmt.Println("GET Employee Called")

	params := mux.Vars(r)
	id := params["id"]

	validID, err := uuid.Parse(id)

	// 400 bad request, invalid uuid
	if err != nil {
		store.SendError(w, http.StatusBadRequest, "GET: Not a valid uuid!", err)
		return
	}

	//Get your employeeItem out of DynamoDb
	output, getEmployeeErr := s.Db.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"EmployeeId": &dynamodb.AttributeValue{
				S: aws.String(validID.String()),
			},
		},
		TableName: aws.String(TableName),
	})

	//can't successfully return output from getItem
	if getEmployeeErr != nil {
		store.SendError(w, http.StatusInternalServerError, "GET: Unable to retrieve that employee from DB!", getEmployeeErr)
		return
	}

	//Check for employee existence in Dynamo
	//If dynamo key doesn't exist it will return a status code of 200, with no data
	if len(output.Item) == 0 {
		store.SendError(w, http.StatusNotFound, "GET: Unable to find employee with that UUID!", errors.New("Output.Item has no length"))
		return
	}

	employeeInfo := store.Employee{}
	unmarshallErr := dynamodbattribute.UnmarshalMap(output.Item, &employeeInfo)
	fmt.Println(" After Unmarshal MAP: STORE.EMPLOYEE: ", output.Item)

	//check for unmarshall errors when unmarshall fields out of output.Item into Employee struct
	// If not nil there's a problem unmarshalling your employee from Dynamo
	if unmarshallErr != nil {
		store.SendError(w, http.StatusInternalServerError, "GET: Employee ID, not found-- unmarshall-able!", unmarshallErr)
		return
	}

	// 200 success, uuid found & item pulled out of dynamoDB
	w.WriteHeader(http.StatusOK)
	resp := GetEmployeeResponse{Status: 200, Data: employeeInfo}
	json.NewEncoder(w).Encode(resp)
	fmt.Println("Get Employee Successfully Completed")
}
