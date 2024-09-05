package user

type LoginRequest struct {
	Email    string `mapstructure:"email" json:"email"`
	Password string `mapstructure:"password" json:"password"`
}

type CreateUserRequest struct {
	Email    string `mapstructure:"email" json:"email"`
	Password string `mapstructure:"password" json:"password"`
}
