package store

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"net/http"
)

//may want to move this error struct

type Error struct {
	Status int    `json:"status,omitempty"`
	Msg    string `json:"msg,omitempty"`
}

type Employee struct {
	Name       string `json:"employeeName,omitempty"`
	Email      string `json:"email,omitempty"`
	EmployeeId string `json:"EmployeeId,omitempty"`
	City       string `json:"city,omitempty"`
	Address    string `json:"address,omitempty"`
	Department string `json:"department,omitempty"`
}

type SubsetDynamoDb interface {
	DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

func SendError(w http.ResponseWriter, statusCode int, errMsg string, err error) {

	//Send back a string with the entered message, plus show the error.
	log := fmt.Sprintf("%v, %v ", errMsg, err.Error())
	fmt.Println(log)

	errResp := Error{statusCode, errMsg}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(errResp); err != nil {
		fmt.Println("Send Error Request failed to Send!")
	}
}
