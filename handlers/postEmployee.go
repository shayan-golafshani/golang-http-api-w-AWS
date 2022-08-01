package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"
)

type GetOneEmployee struct {
	Status   int            `json:"status,omitempty"`
	Employee store.Employee `json:"employee,omitempty"`
}

func (s Server) PostEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST Employee called")

	var newEmployeeUUID = uuid.New()
	fmt.Println("Employee UUID:", newEmployeeUUID)

	//decode the employee from req
	var post store.Employee

	//check you can decode the request
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		store.SendError(w, http.StatusInternalServerError, "POST: Unable to decode request body!", err)
		return
	}
	//check for employee name
	if post.Name == "" {
		store.SendError(w, http.StatusBadRequest, "POST: Reform your request with an employeeName!",
			errors.New("post: Please enter an Employee Name, it's empty"))
		return
	}

	//setup employee to be created
	newEmployee := store.Employee{
		Name:       post.Name,
		Email:      post.Name + "@twilio.com",
		EmployeeId: newEmployeeUUID.String(),
		City:       post.City,
		Address:    post.Address,
		Department: post.Department,
	}
	fmt.Println("New employee Info:", newEmployee)

	//store newEmployee to DynamoDB after marshalMapping them
	dynamoDbEmployee, errMarshalling := dynamodbattribute.MarshalMap(newEmployee)
	if errMarshalling != nil {
		store.SendError(w, http.StatusInternalServerError,
			"POST: Unable to MarshallMap your employee! Try again", errMarshalling)
		return
	}

	//setup putItemInput to create employee record in DynamoDb
	setupItem := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      dynamoDbEmployee,
	}

	if _, errPutItem := s.Db.PutItem(setupItem); errPutItem != nil {
		store.SendError(w, http.StatusInternalServerError,
			"POST: failed to put new Employee Record!", errPutItem)
		return
	}

	//POST response created & sent!
	w.WriteHeader(http.StatusCreated)
	resp := GetOneEmployee{http.StatusCreated, newEmployee}
	json.NewEncoder(w).Encode(resp)
	fmt.Println("Post Employee Successfully Completed")
}
