package store

import "github.com/aws/aws-sdk-go/service/dynamodb"

//WILL MOVE ERROR SOMEWHERE ELSE....
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
