package auth

type SignUpRequest struct {
	Email         string `json:"email" validate:"required,email" example:"user@mail.com"`
	Password      string `json:"password" validate:"required,min=5,max=50" example:"supersecretpassword"`
	PasswordAgain string `json:"password_again" validate:"required,eqfield=Password"  example:"supersecretpassword"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"  example:"user@mail.com"`
	Password string `json:"password" validate:"required,min=5,max=50"  example:"supersecretpassword"`
}

type SignInResponse struct {
	Email string `json:"email"  example:"user@mail.com"`
}
