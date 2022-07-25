package handlers

import (
	"encoding/json"
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
	fmt.Println("POST EMPLOYEE CALLED")

	var newEmployeeUUID = uuid.New()
	fmt.Println("Employee UUID:", newEmployeeUUID)

	//decode the employee from what's posted.
	var post store.Employee

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		fmt.Println("POST: You've got request body issues", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "POST: Unable to decode request body! Try again"})
		return
	}
	if post.Name == "" {
		fmt.Println("POST: Please enter an Employee Name, it's empty.")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "POST: Reform your request with an employeeName!"})
		return
	}

	newEmployee := store.Employee{
		Name:       post.Name,
		Email:      post.Name + "@twilio.com",
		EmployeeId: newEmployeeUUID.String(),
		City:       post.City,
		Address:    post.Address,
		Department: post.Department,
	}
	fmt.Println("New employee:", newEmployee)

	//store newEmployee to DynamoDB after marshalMapping them
	dynamoDbEmployee, errMarshalling := dynamodbattribute.MarshalMap(newEmployee)
	if errMarshalling != nil {
		fmt.Sprintf("POST: failed to DynamoDB marshal Record, %v", errMarshalling.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "POST: Unable to MarshallMap your employee! Try again"})
	}

	_, errPutItem := s.Db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      dynamoDbEmployee,
	})
	if errPutItem != nil {
		fmt.Sprintf("POST: failed to put Record to DynamoDB, %v", errPutItem.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "POST: failed to put Record to DynamoDB!"})
	}

	//send response created!
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := GetOneEmployee{http.StatusCreated, newEmployee}
	json.NewEncoder(w).Encode(resp)
}
