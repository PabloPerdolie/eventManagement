package domain

type UserUpdateRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email" binding:"omitempty,email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UserIdRequest struct {
	UserId int `uri:"id" binding:"required"`
}

type UserListResponse struct {
	Users   []UserResponse `json:"users"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}

type UserRoleUpdateRequest struct {
	Role UserRole `json:"role" binding:"required"`
}

type UserStatusUpdateRequest struct {
	IsActive bool `json:"is_active"`
}
