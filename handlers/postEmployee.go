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

	var newEmployeeUUID = uuid.New()
	fmt.Println("Employee UUID:", newEmployeeUUID)

	//decode the employee from what's posted.
	var post store.Employee

	err := json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		fmt.Println("You've got request body issues", err)
	}
	if post.Name == "" {
		fmt.Println("Please enter an Employee Name")
	}

	newEmployee := store.Employee{
		Name:       post.Name,
		Email:      (post.Name + "@twilio.com"),
		EmployeeId: newEmployeeUUID.String(),
		City:       post.City,
		Address:    post.Address,
		Department: post.Department,
	}

	fmt.Println("New recruit ->", newEmployee)

	//store the employee data locally
	//store.Employees[newEmployee.EmployeeId] = newEmployee

	//store this to DynamoDB
	dynamoDbEmployee, err := dynamodbattribute.MarshalMap(newEmployee)
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}

	_, err = s.Db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      dynamoDbEmployee,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to put Record to DynamoDB, %v", err))
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	resp := GetOneEmployee{http.StatusCreated, newEmployee}

	json.NewEncoder(w).Encode(resp)
}
