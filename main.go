package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/handlers"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/store"

	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
)

const (
	sgGatewayJWTKey        string = "jwt"
	sgGatewayResellerIDKey string = "resellerid"
	sgGatewayScopesKey     string = "scopes"
	sgGatewayUserIDKey     string = "userid"

	contentTypeHeader          = "Content-Type"
	contentTypeApplicationJSON = "application/json"
)

// GetRouter builds the mux router used by a server
func GetRouter(server *Server) *mux.Router {
	fmt.Println("getting routert on Lambda")
	//ADD EMPLOYEES TO THE STORE... will need to store these bad bois in DynamoDB for persistence..
	store.AddEmployees()

	//New mux.Router
	r := mux.NewRouter()

	//Create employee (C)
	r.HandleFunc("/v1/employee/add", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostEmployee(w, r)
	}).Methods("POST")

	//Get Employee via employeeID (R)
	r.HandleFunc("/v1/employee/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetEmployee(w, r)
	}).Methods("GET")

	//Update employee (U)
	r.HandleFunc("/v1/employee/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateEmployee(w, r)
	}).Methods("PATCH")

	//Delete employee (D)
	r.HandleFunc("/v1/employee/delete", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteEmployee(w, r)
	}).Methods("DELETE")

	return r
}

func main() {
	handler, err := Router()
	if err != nil {
		fmt.Println("Unable to create http router for lambda:httpServer", err)
	}
	awslambda.Start(Handle(handler.router))
}

//MCAUTO BOILER PLATE BELOW...
//IGNORE FOR NOW!
type Server struct {
	router http.Handler
}

// Router builds the mux router for the app and returns a Server
func Router(opts ...func(*Server)) (*Server, error) {
	server, err := NewServer(opts...)
	if err != nil {
		return nil, err
	}

	r := GetRouter(server)

	server.router = r
	return server, nil
}

func NewServer(opts ...func(*Server)) (*Server, error) {
	server := &Server{}

	for _, opt := range opts {
		opt(server)
	}
	return server, nil
}

func Handle(router http.Handler) func(ae events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ae events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		req, err := GetRequest(context.Background(), ae)
		if err != nil {
			log.Printf("could not generate request for handlers: %s\n", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		resp, err := getJSONResponse(w)
		if err != nil {
			log.Printf("error getting response: %s\n", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		return resp, nil
	}
}

// GetRequest constructs an http.Request from an events.APIGatewayProxyRequest,
func GetRequest(ctx context.Context, ae events.APIGatewayProxyRequest) (*http.Request, error) {
	var body io.Reader
	if ae.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(ae.Body))
	} else {
		body = bytes.NewBufferString(ae.Body)
	}

	req, err := http.NewRequestWithContext(ctx, ae.HTTPMethod, ae.Path, body)
	if err != nil {
		return nil, err
	}

	// add in query params for things like limit/offset
	query := req.URL.Query()
	for k, v := range ae.QueryStringParameters {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	// add headers for things like x-mako
	for h, v := range ae.Headers {
		req.Header.Set(h, v)
	}

	return req.WithContext(ctx), nil
}

func getJSONResponse(w *httptest.ResponseRecorder) (events.APIGatewayProxyResponse, error) {
	headers := getHeaders(w)
	// enforce application/json response header
	headers[contentTypeHeader] = contentTypeApplicationJSON

	return buildResponse(w, headers)
}

func getResponse(w *httptest.ResponseRecorder) (events.APIGatewayProxyResponse, error) {
	headers := getHeaders(w)

	return buildResponse(w, headers)
}

func getHeaders(w *httptest.ResponseRecorder) (headers map[string]string) {
	headers = make(map[string]string, len(w.HeaderMap))
	for k := range w.HeaderMap {
		headers[k] = w.HeaderMap.Get(k)
	}

	return headers
}

func buildResponse(w *httptest.ResponseRecorder, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var resp events.APIGatewayProxyResponse
	res := w.Result()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(res.Body)
	if err != nil {
		return resp, err
	}
	defer res.Body.Close()

	resp = events.APIGatewayProxyResponse{
		StatusCode: res.StatusCode,
		Body:       buf.String(),
		Headers:    headers,
	}

	return resp, nil
}
