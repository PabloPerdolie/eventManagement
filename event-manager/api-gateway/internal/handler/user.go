package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/event-management/api-gateway/internal/domain"
)

// GetCurrentUser returns the currently authenticated user
// @Summary Get current user
// @Description Returns the profile of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} domain.UserResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /users/me [get]
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Not authenticated",
		})
		return
	}

	user, err := h.service.User.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates the current user's profile
// @Summary Update profile
// @Description Update the profile of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body domain.UserUpdateRequest true "User details to update"
// @Success 200 {object} domain.UserResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /users/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Not authenticated",
		})
		return
	}

	var req domain.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad request",
			Message: err.Error(),
		})
		return
	}

	user, err := h.service.User.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword changes the current user's password
// @Summary Change password
// @Description Change the password of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body domain.ChangePasswordRequest true "Old and new passwords"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /users/me/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Not authenticated",
		})
		return
	}

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad request",
			Message: err.Error(),
		})
		return
	}

	err = h.service.Password.ChangePassword(c.Request.Context(), userID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Password successfully changed",
	})
}

// DeleteAccount deactivates the current user's account
// @Summary Delete account
// @Description Delete (deactivate) the currently authenticated user's account
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} domain.SuccessResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /users/me [delete]
func (h *Handler) DeleteAccount(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Not authenticated",
		})
		return
	}

	err = h.service.User.DeactivateUser(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Account successfully deactivated",
	})
}

// GetAllUsers returns a list of all users (admin only)
// @Summary Get all users
// @Description Get a list of all users (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Page size (default: 10)"
// @Success 200 {object} domain.UserListResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /admin/users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	// Pagination parameters
	page := getIntQueryParam(c, "page", 1)
	limit := getIntQueryParam(c, "limit", 10)
	
	// Calculate offset
	offset := (page - 1) * limit

	users, total, err := h.service.User.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		h.handleError(c, err)
		return
	}

	domainUsers := make([]domain.UserResponse, len(users))
	for i, user := range users {
		domainUsers[i] = domain.UserResponse{
			ID:              user.ID,
			Email:           user.Email,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Role:            domain.UserRole(user.Role),
			IsEmailVerified: user.IsEmailVerified,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, domain.UserListResponse{
		Users:   domainUsers,
		Total:   total,
		Page:    page,
		PerPage: limit,
	})
}

// GetUserById returns a user by ID (admin only)
// @Summary Get user by ID
// @Description Get a user by their ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} domain.UserResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /admin/users/{id} [get]
func (h *Handler) GetUserById(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad request",
			Message: "Invalid user ID",
		})
		return
	}

	user, err := h.service.User.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates a user by ID (admin only)
// @Summary Update user by ID
// @Description Update a user by their ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param request body domain.UserUpdateRequest true "User details to update"
// @Success 200 {object} domain.UserResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /admin/users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad request",
			Message: "Invalid user ID",
		})
		return
	}

	var req domain.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad request",
			Message: err.Error(),
		})
		return
	}

	user, err := h.service.User.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser completely removes a user by ID (admin only)
// @Summary Delete user by ID
// @Description Permanently delete a user by their ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /admin/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad request",
			Message: "Invalid user ID",
		})
		return
	}

	err = h.service.User.HardDeleteUser(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "User successfully deleted",
	})
}
