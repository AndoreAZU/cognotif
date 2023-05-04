package request

type RegisterCustomer struct {
	Name     string `json:"name" validate:"required,alphawithspace"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
