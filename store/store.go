package store

//WILL MOVE ERROR SOMEWHERE ELSE....
type Error struct {
	Status int    `json:"status,omitempty"`
	Msg    string `json:"msg,omitempty"`
}
type Employee struct {
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
	EmployeeId string `json:"EmployeeId,omitempty"`
	City       string `json:"city,omitempty"`
	Address    string `json:"address,omitempty"`
	Department string `json:"department,omitempty"`
}
