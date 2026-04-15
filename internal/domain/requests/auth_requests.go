package requests

type LoginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
