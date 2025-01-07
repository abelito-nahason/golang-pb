package auth

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RoleRoute struct {
	Accessible bool `json:"accessible"`
}
