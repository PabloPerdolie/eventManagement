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

// Create
// @Summary Create a new comment
// @Description Create a new comment
// @Tags comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body domain.CommentCreateRequest true "Comment data"
// @Success 201 {object} domain.AuthResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 409 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /comments/create [post]
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
