package store

import (
	"github.com/google/uuid"
)

type Employee struct {
	Name       string    `json:"name,omitempty"`
	Email      string    `json:"email,omitempty"`
	EmployeeId uuid.UUID `json:"employeeId,omitempty"`
	City       string    `json:"city,omitempty"`
	Address    string    `json:"address,omitempty"`
	Department string    `json:"department,omitempty"`
}

//map to hold all the employees...
var Employees = make(map[uuid.UUID]Employee)

//UUID 1
var uuid1, _ = uuid.Parse("28cdb345-63fc-459f-a1ba-b588630a9fc0")

//UUID 2
var uuid2, _ = uuid.Parse("663c9932-383d-44d1-8b9d-1071731a6312")

var Employee1 = Employee{
	Name:       "Tim Horton",
	Email:      "tim_horton@twilio.com",
	EmployeeId: uuid1,
	City:       "Los Angeles",
	Address:    "222 Boohoo Street",
	Department: "Engineering",
}

var Employee2 = Employee{
	Name:       "Ash Ketchum",
	Email:      "ash_ketchum@twilio.com",
	EmployeeId: uuid2,
	City:       "Detroit",
	Address:    "999 Whine Street",
	Department: "HR",
}

Employees[uuid1] = Employee1
Employees[uuid2] = Employee2


fmt.Println(Employees)