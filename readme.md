# golang-http-api-w-AWS
Setting up an employee CRUD HTTP server in GO that makes use of DynamoDB, Lambda, &amp; API-Gateway.

Please make requests to base url. ->
BASE URL: https://8rj31kamv4.execute-api.us-west-2.amazonaws.com/Dev

CRUD METHODS:

# POST employee
## Endpoint
/v1/employee/add

### Example JSON body in req:
{
"employeeName": "name",
"city": "Denver",
"address": "1 Xander Ave",
"department": "Eng"
}

### Example returned json response:
{"status":201,"employee":{"employeeName":"name","email":"name@twilio.com","EmployeeId":"0ea9875d-13bb-4dd6-83ce-b130b983e905","city":"Denver","address":"1 Xander Ave","department":"Eng"}}

# GET employee
## Endpoint
/v1/employee/{userId}

### Example returned json response:
{"status":200,"data":{"employeeName":"name","email":"name@twilio.com","EmployeeId":"0ea9875d-13bb-4dd6-83ce-b130b983e905","city":"Denver","address":"1 Xander Ave","department":"Eng"}}

# PATCH employee
## Endpoint
/v1/employee/{userId}/update

Can update any of the fields except employeeId
### Example JSON body in req:
{
"department": "HR"
}

### Example returned json response:
{"status":200,"employee":{"employeeName":"name","email":"name@twilio.com","EmployeeId":"0ea9875d-13bb-4dd6-83ce-b130b983e905","city":"Denver","address":"1 Xander Ave","department":"HR"}}

# DELETE employee

## Endpoint
/v1/employee/delete

### Example JSON body in req:
{
"employeeId": "65e69206-316a-4027-bf1a-a116fafdcf54"
}

### Example returned json response:
{
"status":204
}
