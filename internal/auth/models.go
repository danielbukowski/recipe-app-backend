package auth

type SignUpRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	PasswordAgain string `json:"password_again"`
}
