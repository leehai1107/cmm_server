package request

// Login represents user login credentials
// @Description User login request
type Login struct {
	// User's email address
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	// User's password
	Password string `json:"password" binding:"required" example:"password123"`
}

// Register represents user registration data
// @Description User registration request
type Register struct {
	// User's first name
	FirstName string `json:"first_name" binding:"required" example:"John"`
	// User's last name
	LastName string `json:"last_name" binding:"required" example:"Doe"`
	// User's email address
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	// User's password
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}
