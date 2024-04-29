package models
import(

)
type User struct {
	Username string `json:"username" validate:"required,min=5,max=20"`
	Password string `json:"password" validate:"min=6,max=10"`
	Email    string `json:"email" validate:"required,email"`
}
