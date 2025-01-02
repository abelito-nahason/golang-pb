package auth

type RegisterInput struct {
	Email    string `json:email`
	Password string `json:password`
	Role     string `json:role`
}

type LoginInput struct {
	Email    string `json:email`
	Password string `json:password`
}
