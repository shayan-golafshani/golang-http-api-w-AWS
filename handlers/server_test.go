package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sendgrid/mcauto/_vendor/github.com/google/uuid"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

type MockDynamoDB struct {
	GetOutput dynamodb.GetItemOutput
	GetError  error

	UpdateOutput dynamodb.UpdateItemOutput
	UpdateError  error

	PatchOutput dynamodb.PutItemOutput
	PatchError  error

	DeleteOutput dynamodb.DeleteItemOutput
	DeleteError  error
}

//made these pointers so that I can change the value with the mock through the handler function!
//https://stackoverflow.com/questions/27775376/value-receiver-vs-pointer-receiver
func (m *MockDynamoDB) GetItem(_ *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &m.GetOutput, m.GetError
}

func (m *MockDynamoDB) UpdateItem(_ *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return &m.UpdateOutput, m.UpdateError
}

func (m *MockDynamoDB) DeleteItem(_ *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return &m.DeleteOutput, m.DeleteError
}

func (m *MockDynamoDB) PutItem(_ *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &m.PatchOutput, m.PatchError
}

func TestServer_GetEmployee(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server handler")
}

var _ = Describe("Employee CRUD Server tests", func() {
	var (
		server Server
		m      MockDynamoDB
	)

	BeforeEach(func() {
		m = MockDynamoDB{}
		server = Server{}

		server.Router = GetRouter(&server)
		server.Db = &m
	})

	Context("Happy path tests", func() {
		It("Get Employee handler Success", func() {

			req, err := http.NewRequest("GET", "/v1/employee/{id}", nil)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some fake UUID
			vars := map[string]string{"id": "663c9932-383d-44d1-8b9d-1071731a6312"}

			req = mux.SetURLVars(req, vars)

			newEmployee := store.Employee{
				Name:       "shayan",
				Email:      "shayan@twilio.com",
				EmployeeId: "1",
				City:       "Yar",
				Address:    "The Black Pearl",
				Department: "Piratology",
			}

			item, _ := dynamodbattribute.MarshalMap(newEmployee)

			m.GetError = nil
			m.GetOutput = dynamodb.GetItemOutput{
				Item: item,
			}

			server.GetEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp GetEmployeeResponse

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusOK))
			Expect(resp.Data).To(Equal(newEmployee))
		})
		It("Post Employee handler Success", func() {
			reader := strings.NewReader(`{"employeeName": "pablo", "city":  "Denver", "address": "3 Main St", "department": "eng"}`)

			req, err := http.NewRequest("POST", "/v1/employee/add", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			server.PostEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp GetOneEmployee

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusCreated))
			Expect(resp.Employee.Name).To(Equal("pablo"))
			Expect(resp.Employee.Email).To(Equal("pablo@twilio.com"))
			//test for UUID REGexp
			r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
			val := r.MatchString(resp.Employee.EmployeeId)
			Expect(val).To(Equal(true))
			Expect(resp.Employee.City).To(Equal("Denver"))
			Expect(resp.Employee.Address).To(Equal("3 Main St"))
			Expect(resp.Employee.Department).To(Equal("eng"))
		})
		It("Delete Employee handler Success", func() {
			//some UUID
			uuid := uuid.New().String()
			reader := strings.NewReader(fmt.Sprintf(`{"employeeId": "%v"}`, uuid))

			req, err := http.NewRequest("DELETE", "/v1/employee/delete", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			server.DeleteEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp DeletionResponse

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusNoContent))
			Expect(resp.Msg).To(Equal("Successful deletion of Employee"))
		})
		It("Update Employee handler Success", func() {
			reader := strings.NewReader(`{"employeeName": "pablo", "email":"pablo@yahoo.com", "city":  "Denver", "address": "3 Main St", "department": "eng"}`)

			req, err := http.NewRequest("POST", "/v1/employee/{id}/update", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some fake uuid
			vars := map[string]string{"id": "663c9932-383d-44d1-8b9d-1071731a6312"}

			req = mux.SetURLVars(req, vars)

			//newEmployee := store.Employee{
			//	Name:       "pablo",
			//	Email:      "pablo@yahoo.com",
			//	EmployeeId: "663c9932-383d-44d1-8b9d-1071731a6312",
			//	City:       "Denver",
			//	Address:    "3 Main St",
			//	Department: "eng",
			//}

			//item, _ := dynamodbattribute.MarshalMap(newEmployee)

			//update error
			//m.UpdateError = nil
			//m.UpdateOutput = dynamodb.UpdateItemOutput{
			//	Attributes: item,
			//}

			server.UpdateEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp GetOneEmployee

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusOK))
			//Expect(resp.Employee.Name).To(Equal("pablo"))
			//Expect(resp.Employee.Email).To(Equal("pablo@twilio.com"))
			//test for UUID REGexp
			//r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
			//val := r.MatchString(resp.Employee.EmployeeId)
			//Expect(val).To(Equal(true))
			//Expect(resp.Employee.City).To(Equal("Denver"))
			//Expect(resp.Employee.Address).To(Equal("3 Main St"))
			//Expect(resp.Employee.Department).To(Equal("eng"))
		})
	})

	Context("Sad Path tests", func() {
		It("Should test Invalid UUID passed to GET Employee Handler", func() {

			req, err := http.NewRequest("GET", "/v1/employee/{id}", nil)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some fake UUID
			vars := map[string]string{"id": "ZZ3c9932-383d-44d1-8b9d-1071731a6312"}

			req = mux.SetURLVars(req, vars)

			server.GetEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp store.Error

			err = json.NewDecoder(res.Body).Decode(&resp)

			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			//Update tests
			Expect(resp.Status).To(Equal(http.StatusBadRequest))
			Expect(resp.Msg).To(Equal("GET: Not a valid uuid!"))

		})
		It("Should test GET EMPLOYEE, getItem method in fake DynamoDB client", func() {

			req, err := http.NewRequest("GET", "/v1/employee/{id}", nil)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some fake UUID
			vars := map[string]string{"id": "663c9932-383d-44d1-8b9d-1071731a6312"}

			req = mux.SetURLVars(req, vars)

			m.GetError = errors.New("error getting EmployeeID out of DynamoDB")
			m.GetOutput = dynamodb.GetItemOutput{}

			server.GetEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp store.Error

			err = json.NewDecoder(res.Body).Decode(&resp)

			Expect(resp.Status).To(Equal(http.StatusInternalServerError))
			Expect(resp.Msg).To(Equal("GET: Unable to find employee! Try again"))
		})
		It("Should test GET EMPLOYEE, getItem method for missing not found user", func() {

			req, err := http.NewRequest("GET", "/v1/employee/{id}", nil)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some fake UUID
			vars := map[string]string{"id": "663c9932-383d-44d1-8b9d-1071731a6312"}

			req = mux.SetURLVars(req, vars)

			m.GetError = nil
			m.GetOutput = dynamodb.GetItemOutput{}

			server.GetEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp store.Error

			err = json.NewDecoder(res.Body).Decode(&resp)

			Expect(resp.Status).To(Equal(http.StatusNotFound))
			Expect(resp.Msg).To(Equal("GET: Unable to find employee! Try again"))
		})
		It("Should test invalid  POST response body issues, for example improper json", func() {
			//missing comma after address
			reader := strings.NewReader(`{"employeeName": "pablo", "city":  "Denver", "address": "3 Main St" "department": "eng"}`)

			req, err := http.NewRequest("POST", "/v1/employee/add", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			server.PostEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp store.Error

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusInternalServerError))
			Expect(resp.Msg).To(Equal("POST: Unable to decode request body! Try again"))
		})
		It("Should test for an empty name in POST body issues", func() {
			//missing "employeeName": "pablo", from beginning of JSON response
			reader := strings.NewReader(`{"city":  "Denver", "address": "3 Main St", "department": "eng"}`)

			req, err := http.NewRequest("POST", "/v1/employee/add", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			server.PostEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp store.Error

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusBadRequest))
			Expect(resp.Msg).To(Equal("POST: Reform your request with an employeeName!"))
		})
		It("Should test invalid PUT item error in POST", func() {
			reader := strings.NewReader(`{"employeeName": "pablo", "city":  "Denver", "address": "3 Main St", "department": "eng"}`)

			req, err := http.NewRequest("POST", "/v1/employee/add", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			m.PatchError = errors.New("POST: failed to put Record to DynamoDB")
			m.PatchOutput = dynamodb.PutItemOutput{}

			server.PostEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp store.Error

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusInternalServerError))
			Expect(resp.Msg).To(Equal("POST: failed to put Record to DynamoDB!"))
		})
		It("Should test DELETE Employee request body issues, for example improper json", func() {
			//some UUID
			uuid := uuid.New().String()
			//missing quote after employeeId to break JSON
			reader := strings.NewReader(fmt.Sprintf(`{"employeeId: "%v}`, uuid))

			req, err := http.NewRequest("DELETE", "/v1/employee/delete", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			server.DeleteEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp DeletionResponse

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusBadRequest))
			Expect(resp.Msg).To(Equal("DELETE: Deletion request body config issues."))
		})
		//It("Should test DELETE Employee request body has proper UUID", func() {
		//
		//	//invalid UUID
		//	reader := strings.NewReader(`{"employeeId": "1d7aafd0-5b1f-4eec-8127-639d8ceb5ed0zz"}`)
		//
		//	req, err := http.NewRequest("DELETE", "/v1/employee/delete", reader)
		//	if err != nil {
		//		fmt.Printf("could not create a request: %v", err.Error())
		//	}
		//
		//	w := httptest.NewRecorder()
		//
		//	server.DeleteEmployee(w, req)
		//
		//	res := w.Result()
		//	defer res.Body.Close()
		//
		//	var resp DeletionResponse
		//
		//	err = json.NewDecoder(res.Body).Decode(&resp)
		//	if err != nil {
		//		fmt.Printf("unable to successfully decode response got error %v", err)
		//	}
		//
		//	fmt.Println("Printing resp:", resp)
		//
		//	//refactor
		//	Expect(resp.Status).To(Equal(http.StatusBadRequest))
		//	Expect(resp.Msg).To(Equal("DELETE: Not a valid uuid!"))
		//})

		It("Should test DELETE Employee request body can deal with empty JSON edge-case", func() {

			//empty JSON turns out this case is also taken care of....
			reader := strings.NewReader(`{}`)

			req, err := http.NewRequest("DELETE", "/v1/employee/delete", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			server.DeleteEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp DeletionResponse

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusBadRequest))
			Expect(resp.Msg).To(Equal("DELETE: Cannot pass in an empty body, attach a valid employeeId to delete."))
		})
		It("Should test DELETE invalid Delete Item error", func() {
			//some UUID
			uuid := uuid.New().String()
			reader := strings.NewReader(fmt.Sprintf(`{"employeeId": "%v"}`, uuid))

			req, err := http.NewRequest("DELETE", "/v1/employee/delete", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			m.DeleteError = errors.New("DELETE: failed to delete record from DynamoDB")
			m.DeleteOutput = dynamodb.DeleteItemOutput{}

			server.DeleteEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp DeletionResponse

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusInternalServerError))
			Expect(resp.Msg).To(Equal("DELETE: failed to delete Record from DynamoDB!"))
		})
		It("Should test Update unknown fields in JSON", func() {
			//extra field
			reader := strings.NewReader(`{"TEEEHEE": "UNKOWN!!", "employeeName": "pablo", "email":"pablo@yahoo.com", "city":  "Denver", "address": "3 Main St", "department": "eng"}`)

			req, err := http.NewRequest("POST", "/v1/employee/{id}/update", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some fake uuid
			vars := map[string]string{"id": "663c9932-383d-44d1-8b9d-1071731a6312"}

			req = mux.SetURLVars(req, vars)

			server.UpdateEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp store.Error

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusBadRequest))
			Expect(resp.Msg).To(Equal("PATCH: Unknown field(s) included in request body or empty request body. Please only use editable employee information."))
		})
		It("Should test Update invalid UUID field in JSON", func() {
			reader := strings.NewReader(`{"employeeName": "pablo", "email":"pablo@yahoo.com", "city":  "Denver", "address": "3 Main St", "department": "eng"}`)

			req, err := http.NewRequest("POST", "/v1/employee/{id}/update", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some fake invalid uuid
			vars := map[string]string{"id": "ZZ3c9932-383d-44d1-8b9d-1071731a6312"}

			req = mux.SetURLVars(req, vars)

			server.UpdateEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp store.Error

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusBadRequest))
			Expect(resp.Msg).To(Equal("PATCH: not a valid uuid!"))
		})
		It("Should test Update's GET-ITEM method", func() {
			reader := strings.NewReader(`{"employeeName": "pablo", "email":"pablo@yahoo.com", "city":  "Denver", "address": "3 Main St", "department": "eng"}`)

			req, err := http.NewRequest("POST", "/v1/employee/{id}/update", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some fake invalid uuid
			vars := map[string]string{"id": "663c9932-383d-44d1-8b9d-1071731a6312"}

			req = mux.SetURLVars(req, vars)

			m.GetError = errors.New("patch: Unable to find employee for GetItem")
			m.GetOutput = dynamodb.GetItemOutput{}

			server.UpdateEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp store.Error

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			fmt.Println("Printing resp:", resp)

			//refactor
			Expect(resp.Status).To(Equal(http.StatusInternalServerError))
			Expect(resp.Msg).To(Equal("PATCH: Unable to find employee!"))
		})
		It("Should test Update's Update-ITEM method", func() {
			reader := strings.NewReader(`{"employeeName": "pablo", "email":"pablo@yahoo.com", "city":  "Denver", "address": "3 Main St", "department": "eng"}`)

			req, err := http.NewRequest("POST", "/v1/employee/{id}/update", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some fake invalid uuid
			vars := map[string]string{"id": "663c9932-383d-44d1-8b9d-1071731a6312"}

			req = mux.SetURLVars(req, vars)

			m.UpdateError = errors.New("patch: Unable to find employee for GetItem")
			m.UpdateOutput = dynamodb.UpdateItemOutput{}

			server.UpdateEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp store.Error

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			Expect(resp.Status).To(Equal(http.StatusInternalServerError))
			Expect(resp.Msg).To(Equal("PATCH: FAILED TO UPDATE YOUR EMPLOYEE RECORD!"))
		})
	})

})
