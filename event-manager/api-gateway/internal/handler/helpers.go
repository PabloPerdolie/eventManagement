package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/event-management/api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var errorStatusMap = map[error]int{
	domain.ErrInvalidCredentials:   http.StatusUnauthorized,
	domain.ErrEmailAlreadyExists:    http.StatusConflict,
	domain.ErrUserNotFound:         http.StatusNotFound,
	domain.ErrInvalidToken:         http.StatusUnauthorized,
	domain.ErrTokenExpired:         http.StatusUnauthorized,
	domain.ErrInvalidResetToken:    http.StatusBadRequest,
	domain.ErrUserNotActive:        http.StatusForbidden,
	//domain.ErrInsufficientPermission: http.StatusForbidden,
}

func (h *Handler) handleError(c *gin.Context, err error) {
	var statusCode int
	var errorResponse domain.ErrorResponse

	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		statusCode = http.StatusUnauthorized
		errorResponse = domain.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid credentials",
		}
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		statusCode = http.StatusConflict
		errorResponse = domain.ErrorResponse{
			Error:   "Conflict",
			Message: "User already exists",
		}
	case errors.Is(err, domain.ErrUserNotFound):
		statusCode = http.StatusNotFound
		errorResponse = domain.ErrorResponse{
			Error:   "NotFound",
			Message: "User not found",
		}
	case errors.Is(err, domain.ErrInvalidToken):
		statusCode = http.StatusUnauthorized
		errorResponse = domain.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid token",
		}
	case errors.Is(err, domain.ErrTokenExpired):
		statusCode = http.StatusUnauthorized
		errorResponse = domain.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Token has expired",
		}
	case errors.Is(err, domain.ErrInvalidResetToken):
		statusCode = http.StatusBadRequest
		errorResponse = domain.ErrorResponse{
			Error:   "BadRequest",
			Message: "Invalid reset token",
		}
	case errors.Is(err, domain.ErrUserNotActive):
		statusCode = http.StatusForbidden
		errorResponse = domain.ErrorResponse{
			Error:   "Forbidden",
			Message: "User is not active",
		}
	default:
		statusCode = http.StatusInternalServerError
		errorResponse = domain.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		}
	}

	c.JSON(statusCode, errorResponse)
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID is not of type UUID")
	}

	return id, nil
}

func getIntQueryParam(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.DefaultQuery(key, strconv.Itoa(defaultValue))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
