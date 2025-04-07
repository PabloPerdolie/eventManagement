package handler

import (
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
)

type CommentService interface {
	Create(ctx context.Context, request domain.CommentCreateRequest) error
}

type Comment struct {
	commentService CommentService
}

func NewComment(commentService CommentService) Comment {
	return Comment{
		commentService: commentService,
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
func (h *Comment) Create(c *gin.Context) {
	var req domain.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad request",
			Message: err.Error(),
		})
		return
	}

	userId, _ := getUserIdFromContext(c)

	req.SenderId = userId

	err := h.commentService.Create(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, "comment created")
}
