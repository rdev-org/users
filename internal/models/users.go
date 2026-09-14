package models

type User struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile,omitempty"`
	Password   string `json:"password,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
}
