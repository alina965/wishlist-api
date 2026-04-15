package domain

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	ID       int    `json:"id"`
	Token    string `json:"token"`
}
