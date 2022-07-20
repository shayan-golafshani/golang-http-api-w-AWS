package handlers

import (
	"encoding/json"
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

type deletionResponse struct {
	Status int
	Msg    string ``
}

func (s Server) DeleteEmployee(w http.ResponseWriter, r *http.Request) {

	var post deletionReq

	deleteErrDynamo := json.NewDecoder(r.Body).Decode(&post)

	if deleteErrDynamo != nil {
		fmt.Println("Deletion Request Body Config Issues", deleteErrDynamo.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "DELETE: Deletion request body config issues."})
		if err != nil {
			fmt.Println("failed to write error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	if post.EmployeeId.String() == "" {
		fmt.Println("Please Submit EmployeeId for deletion")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "DELETE: Please Submit EmployeeId for deletion"})
		return
	}

	//PARSE BODY
	validID, invalidIDErr := uuid.Parse(post.EmployeeId.String())
	if invalidIDErr != nil {
		fmt.Println("DELETE: Invalid employee UUID, not a valid UUID, try again: ", invalidIDErr.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 400, Msg: "DELETE: Not a valid uuid!"})
		return
	}

	//input to delete item
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"EmployeeId": {
				S: aws.String(validID.String()),
			},
		},
		TableName: aws.String(TableName),
	}

	_, deleteErrDynamo = s.Db.DeleteItem(input)
	if deleteErrDynamo != nil {
		fmt.Println("DELETE: Got error calling DeleteItem", deleteErrDynamo.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Error{Status: 500, Msg: "DELETE: failed to delete Record from DynamoDB!"})
		return
	}

	fmt.Println("Deleted employeeId " + validID.String() + " from table" + TableName)

	//DO i still need these checks?
	////TODO
	//_, ok := store.Employees[validID]
	//
	////TODO
	//if !ok {
	//	fmt.Println("Status 404, Employee ID not found.")
	//
	//	w.WriteHeader(http.StatusNotFound)
	//	w.Header().Set("Content-Type", "application/json")
	//	json.NewEncoder(w).Encode(store.Error{Status: 404, Msg: "Employee Id not valid, can't be deleted"})
	//	return
	//}

	// WITH DYNAMODB PARADIGM JUST TRY TO GO FOR THE STRAIGHT DELETE!
	//I could try to do another GET and see if that's successful, but it may be excessive
	//TODO check code below
	//if ok {
	//	delete(store.Employees, validID)
	//}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	resp := deletionResponse{Status: 204, Msg: "Successful deletion of Employee"}
	json.NewEncoder(w).Encode(resp)
	fmt.Println("Removed employee! さよなら！")
}
