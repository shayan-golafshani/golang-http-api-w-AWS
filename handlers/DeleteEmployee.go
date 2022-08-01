package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"net/http"

	"github.com/google/uuid"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"
)

type deletionReq struct {
	EmployeeId uuid.UUID `json:"EmployeeId,omitempty"`
}

type DeletionResponse struct {
	Status int
	Msg    string
}

func (s Server) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE Employee Called")

	var deleteReq deletionReq

	// Decode your request into deletionReq, otherwise your deletion request body isn't configured properly
	if deleteErrDynamo := json.NewDecoder(r.Body).Decode(&deleteReq); deleteErrDynamo != nil {
		store.SendError(w, http.StatusBadRequest, "DELETE: Deletion request body config issues.", deleteErrDynamo)
	}

	//deals with empty JSON being passed and improperly decoded.
	if deleteReq.EmployeeId.String() == "00000000-0000-0000-0000-000000000000" {
		store.SendError(w, http.StatusBadRequest, "DELETE: Cannot pass in an empty body, attach a valid employeeId to delete.",
			errors.New("empty JSON being passed in delete request, need valid EmployeeId to delete"))
		return
	}

	//inputObject required to delete item from DynamoDB
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"EmployeeId": {
				S: aws.String(deleteReq.EmployeeId.String()),
			},
		},
		TableName: aws.String(TableName),
	}

	if _, deleteItemErrDynamo := s.Db.DeleteItem(input); deleteItemErrDynamo != nil {
		store.SendError(w, http.StatusInternalServerError, "DELETE: failed to delete record", deleteItemErrDynamo)
		return
	}

	fmt.Sprintf("Deleted employee Id %v, from table: %v", deleteReq.EmployeeId.String(), TableName)

	w.WriteHeader(http.StatusNoContent)
	resp := DeletionResponse{Status: 204, Msg: "Successful Deletion of Employee"}
	json.NewEncoder(w).Encode(resp)
	fmt.Println("Delete Employee Successfully Completed")
}
