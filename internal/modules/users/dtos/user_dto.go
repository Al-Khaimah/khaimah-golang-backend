package users

type SignupRequestDTO struct {
	FirstName  string   `json:"first_name" validate:"required"`
	LastName   string   `json:"last_name" validate:"not required"`
	Email      string   `json:"email" validate:"required,email"`
	Categories []string `json:"categories" validate:"required,email"`
	Password   string   `json:"password" validate:"required"`
}

type SignupResponseDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
