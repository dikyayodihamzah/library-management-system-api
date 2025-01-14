package user

type LoginRequest struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

type SignUpRequest struct {
	Name string `json:"name,omitempty" validate:"required"`
	LoginRequest
}

type LoginResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserData     *User  `json:"user_data,omitempty"`
}
