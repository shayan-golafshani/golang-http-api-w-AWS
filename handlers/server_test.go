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
		server.SetupRouter()
		server.Db = &m
	})

	Context("Happy Path tests", func() {
		It("Status:200 Get Employee Success", func() {

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

			Expect(resp.Status).To(Equal(http.StatusOK))
			Expect(resp.Data).To(Equal(newEmployee))
		})
		It("Status:201 Post Employee Success", func() {
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

			Expect(resp.Status).To(Equal(http.StatusCreated))
			Expect(resp.Employee.Name).To(Equal("pablo"))
			Expect(resp.Employee.Email).To(Equal("pablo@twilio.com"))
			//UUID REGexp
			r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
			val := r.MatchString(resp.Employee.EmployeeId)
			Expect(val).To(Equal(true))
			Expect(resp.Employee.City).To(Equal("Denver"))
			Expect(resp.Employee.Address).To(Equal("3 Main St"))
			Expect(resp.Employee.Department).To(Equal("eng"))
		})
		It("Status:204 Delete Employee Success", func() {
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

			Expect(resp.Status).To(Equal(http.StatusNoContent))
			Expect(resp.Msg).To(Equal("Successful Deletion of Employee"))
		})
		It("Status 200: Update Employee Success", func() {
			reader := strings.NewReader(`{"employeeName": "pablo", "email":"pablo@yahoo.com", "city":  "Denver", "address": "3 Main St", "department": "eng"}`)

			req, err := http.NewRequest("PATCH", "/v1/employee/{id}/update", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some fake uuid
			vars := map[string]string{"id": "663c9932-383d-44d1-8b9d-1071731a6312"}

			req = mux.SetURLVars(req, vars)

			newEmployee := store.Employee{
				Name:       "pablo",
				Email:      "pablo@yahoo.com",
				EmployeeId: "663c9932-383d-44d1-8b9d-1071731a6312",
				City:       "Denver",
				Address:    "3 Main St",
				Department: "sales",
			}

			item, _ := dynamodbattribute.MarshalMap(newEmployee)

			m.GetError = nil
			m.GetOutput = dynamodb.GetItemOutput{
				Item: item,
			}

			m.UpdateError = nil

			server.UpdateEmployee(w, req)

			res := w.Result()
			defer res.Body.Close()

			var resp GetOneEmployee

			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				fmt.Printf("unable to successfully decode response got error %v", err)
			}

			//refactor
			Expect(resp.Status).To(Equal(http.StatusOK))
			//update field for test
			newEmployee.Department = "eng"
			Expect(resp.Employee).To(Equal(newEmployee))

		})
	})

	Context("Sad Path tests", func() {
		It("Should test GET Employee (400), Invalid UUID", func() {

			req, err := http.NewRequest("GET", "/v1/employee/{id}", nil)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some invalid UUID
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

			Expect(resp.Status).To(Equal(http.StatusBadRequest))
			Expect(resp.Msg).To(Equal("GET: Not a valid uuid!"))
		})
		It("Should test GET Employee (404), getItem method error for not found user", func() {

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
			Expect(resp.Msg).To(Equal("GET: Unable to find employee with that UUID!"))
		})
		It("Should test GET Employee (500), getItem method error in mock DynamoDB client", func() {

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
			Expect(resp.Msg).To(Equal("GET: Unable to retrieve that employee from DB!"))
		})
		It("Should test POST Employee (500), request body -- invalid json", func() {
			//missing comma after address field
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

			Expect(resp.Status).To(Equal(http.StatusInternalServerError))
			Expect(resp.Msg).To(Equal("POST: Unable to decode request body!"))
		})
		It("Should test POST Employee (500), empty name field in POST json body ", func() {
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

			Expect(resp.Status).To(Equal(http.StatusBadRequest))
			Expect(resp.Msg).To(Equal("POST: Reform your request with an employeeName!"))
		})
		It("Should test POST Employee (500), putItem method error in mock DynamoDB client", func() {
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

			Expect(resp.Status).To(Equal(http.StatusInternalServerError))
			Expect(resp.Msg).To(Equal("POST: failed to put new Employee Record!"))
		})
		It("Should test DELETE Employee (400), request body -- invalid json", func() {
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

			Expect(resp.Status).To(Equal(http.StatusBadRequest))
			Expect(resp.Msg).To(Equal("DELETE: Deletion request body config issues."))
		})
		It("Should test DELETE Employee (400), request body-- empty json edge-case", func() {
			//empty JSON edge-case, also taken care of....
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

			Expect(resp.Status).To(Equal(http.StatusBadRequest))
			Expect(resp.Msg).To(Equal("DELETE: Cannot pass in an empty body, attach a valid employeeId to delete."))
		})
		It("Should test DELETE Employee (500), deleteItem method error in mock DynamoDB client", func() {
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

			Expect(resp.Status).To(Equal(http.StatusInternalServerError))
			Expect(resp.Msg).To(Equal("DELETE: failed to delete record"))
		})
		It("Should test UPDATE Employee (500), unknown additional field in json", func() {
			//extra field
			reader := strings.NewReader(
				`{"TEEEHEE": "UNKNOWN!!", "employeeName": "pablo", "email":"pablo@yahoo.com", "city":  "Denver", "address": "3 Main St", "department": "eng"}`)

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

			Expect(resp.Status).To(Equal(http.StatusBadRequest))
			Expect(resp.Msg).To(Equal("PATCH: Unknown field(s) included in request body or empty request body. Please only use editable employee information."))
		})
		It("Should test UPDATE Employee (500), invalid UUID field in JSON", func() {
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
		It("Should test UPDATE Employee (500), getItem method error in mock DynamoDB client", func() {
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

			Expect(resp.Status).To(Equal(http.StatusInternalServerError))
			Expect(resp.Msg).To(Equal("PATCH: Unable to find employee!"))
		})
		It("Should test UPDATE Employee (500), updateItem method error in mock DynamoDb client", func() {
			reader := strings.NewReader(`{"employeeName": "pablo", "email":"pablo@yahoo.com", "city":  "Denver", "address": "3 Main St", "department": "eng"}`)

			req, err := http.NewRequest("POST", "/v1/employee/{id}/update", reader)
			if err != nil {
				fmt.Printf("could not create a request: %v", err.Error())
			}

			w := httptest.NewRecorder()

			//some fake uuid
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
