package auth

type SignUpRequest struct {
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=5,max=50"`
	PasswordAgain string `json:"password_again" validate:"required,eqfield=Password"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=50"`
}
