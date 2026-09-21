package dto

type ProfileRequestDTO struct {
	UserID       string             `json:"user_id"`
	FirstName    string             `json:"first_name"`
	MiddleName   string             `json:"middle_name"`
	LastName     string             `json:"last_name"`
	MobileNumber string             `json:"mobile_number"`
	DateOfBirth  string             `json:"date_of_birth"`
	Gender       string             `json:"gender"`
	Avatar       string             `json:"avatar"`
	Address      *AddressRequestDTO `json:"address"`

	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AddressRequestDTO struct {
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2,omitempty"`
	City          string `json:"city"`
	StateProvince string `json:"state_province"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
	IsDefault     bool   `json:"is_default"`
	AddressType   string `json:"address_type"`
}
