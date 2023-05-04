package response

import "time"

type GetCustomer struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}
