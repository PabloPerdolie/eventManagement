package handler

import (
	"context"
	"net/http"

	"github.com/PabloPerdolie/event-manager/api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req domain.UserLoginRequest) (*domain.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, tokenString string) (*domain.JWTClaims, error)
	GetUserInfo(ctx context.Context, userId int) (*domain.UserResponse, error)
	CreatePasswordResetToken(ctx context.Context, email string) (string, error)
}

type Auth struct {
	authService AuthService
}

func NewAuth(authService AuthService) Auth {
	return Auth{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.UserRegisterRequest true "User registration data"
// @Success 201 {object} domain.AuthResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 409 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /auth/register [post]
func (h *Auth) Register(c *gin.Context) {
	var req domain.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad request",
			Message: err.Error(),
		})
		return
	}
	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// Login handles user login
// @Summary Login a user
// @Description Login a user with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.UserLoginRequest true "User credentials"
// @Success 200 {object} domain.AuthResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /auth/login [post]
func (h *Auth) Login(c *gin.Context) {
	var req domain.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad request",
			Message: err.Error(),
		})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Logout handles user logout
// @Summary Logout a user
// @Description Logout a user by invalidating their token
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} domain.SuccessResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /auth/logout [post]
func (h *Auth) Logout(c *gin.Context) {
	token := extractTokenFromHeader(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid or missing token",
		})
		return
	}

	err := h.authService.Logout(c.Request.Context(), token)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Successfully logged out",
	})
}

// RefreshToken refreshes an access token using a refresh token
// @Summary Refresh a token
// @Description Refresh an access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} domain.TokenResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /auth/refresh [post]
func (h *Auth) RefreshToken(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad request",
			Message: err.Error(),
		})
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

//// ForgotPassword initiates the password reset process
//// @Summary Request password reset
//// @Description Request a password reset email
//// @Tags auth
//// @Accept json
//// @Produce json
//// @Param request body domain.ForgotPasswordRequest true "Email address"
//// @Success 200 {object} domain.SuccessResponse
//// @Failure 400 {object} domain.ErrorResponse
//// @Failure 500 {object} domain.ErrorResponse
//// @Router /auth/forgot-password [post]
//func (h *Auth) ForgotPassword(c *gin.Context) {
//	var req domain.ForgotPasswordRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
//			Error:   "Bad request",
//			Message: err.Error(),
//		})
//		return
//	}
//
//	err := h.authService.ForgotPassword(c.Request.Context(), req.Email)
//	if err != nil {
//		handleError(c, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, domain.SuccessResponse{
//		Message: "If your email is registered, you will receive a password reset link",
//	})
//}

//// ResetPassword resets a user's password using a token
//// @Summary Reset password
//// @Description Reset a user's password using a valid reset token
//// @Tags auth
//// @Accept json
//// @Produce json
//// @Param request body domain.ResetPasswordRequest true "Reset token and new password"
//// @Success 200 {object} domain.SuccessResponse
//// @Failure 400 {object} domain.ErrorResponse
//// @Failure 401 {object} domain.ErrorResponse
//// @Failure 500 {object} domain.ErrorResponse
//// @Router /auth/reset-password [post]
//func (h *Auth) ResetPassword(c *gin.Context) {
//	var req domain.ResetPasswordRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
//			Error:   "Bad request",
//			Message: err.Error(),
//		})
//		return
//	}
//
//	err := h.authService.ResetPassword(c.Request.Context(), req)
//	if err != nil {
//		handleError(c, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, domain.SuccessResponse{
//		Message: "Password has been successfully reset",
//	})
//}

func extractTokenFromHeader(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
		return bearerToken[7:]
	}
	return ""
}
