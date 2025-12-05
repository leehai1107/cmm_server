package request

// CreateAccount represents admin account creation data
// @Description Admin account creation request
type CreateAccount struct {
	// User's first name
	FirstName string `json:"first_name" binding:"required" example:"John"`
	// User's last name
	LastName string `json:"last_name" binding:"required" example:"Doe"`
	// User's email address
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	// User's password
	Password string `json:"password" binding:"required,min=8" example:"password123"`
	// User's role
	Role string `json:"role" binding:"required" example:"admin"`
}
