package requests

const (
	minPathLength = 8
)

type BasicAuth struct {
	Email    string `json:"email" validate:"required" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"11111111"`
}

type LoginRequest struct {
	BasicAuth
}

type RegisterRequest struct {
	BasicAuth
	FirstName string `json:"first_name" validate:"required" example:"John"`
	LastName  string `json:"last_name" validate:"required" example:"Doe"`
	Role      string `json:"role" validate:"required" example:"admin"`
	Username  string `json:"username" validate:"required" example:"johndoe"`
}

type RefreshRequest struct {
	Token string `json:"token" validate:"required" example:"refresh_token"`
}
