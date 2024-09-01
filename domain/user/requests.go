package user

type LoginRequest struct {
	Email    string `mapstructure:"email" json:"email"`
	Password string `mapstructure:"password" json:"password"`
}
