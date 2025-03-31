package handler

//import (
//	"net/http"
//	"strconv"
//
//	"github.com/PabloPerdolie/event-manager/api-gateway/internal/domain"
//	"github.com/PabloPerdolie/event-manager/api-gateway/internal/service/comment"
//	"github.com/gin-gonic/gin"
//)
//
//// CreateComment creates a new comment via RabbitMQ message queue
//// @Summary Create a new comment
//// @Description Create a new comment for an event
//// @Tags comments
//// @Accept json
//// @Produce json
//// @Security Bearer
//// @Param request body comment.CreateCommentMessage true "Comment details"
//// @Success 202 {object} domain.SuccessResponse
//// @Failure 400 {object} domain.ErrorResponse
//// @Failure 401 {object} domain.ErrorResponse
//// @Failure 500 {object} domain.ErrorResponse
//// @Router /api/v1/comments [post]
//func (h *Handler) CreateComment(c *gin.Context) {
//	userId, err := getUserIdFromContext(c)
//	if err != nil {
//		h.logger.Errorw("Failed to get user Id from context", "error", err)
//		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
//			Error:   "Unauthorized",
//			Message: "You must be logged in to create a comment",
//		})
//		return
//	}
//
//	var request comment.CreateCommentMessage
//	if err := c.ShouldBindJSON(&request); err != nil {
//		h.logger.Errorw("Invalid request payload", "error", err)
//		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
//			Error:   "Bad Request",
//			Message: "Invalid request payload",
//		})
//		return
//	}
//
//	// Override sender Id with the authenticated user's Id
//	request.SenderId = int(userId)
//
//	// Validate event Id
//	if request.EventId <= 0 {
//		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
//			Error:   "Bad Request",
//			Message: "Invalid event Id",
//		})
//		return
//	}
//
//	// Validate content
//	if request.Content == "" {
//		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
//			Error:   "Bad Request",
//			Message: "Comment content cannot be empty",
//		})
//		return
//	}
//
//	// Create the comment asynchronously
//	err = h.service.Comment.CreateComment(c.Request.Context(), request)
//	if err != nil {
//		h.logger.Errorw("Failed to create comment", "error", err, "event_id", request.EventId, "user_id", userId)
//		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
//			Error:   "Internal Server Error",
//			Message: "Failed to create comment, please try again later",
//		})
//		return
//	}
//
//	c.JSON(http.StatusAccepted, domain.SuccessResponse{
//		Message: "Comment is being processed",
//	})
//}
//
//// GetUserComments gets all comments for a specific user
//// @Summary Get comments for user
//// @Description Get all comments for the authenticated user
//// @Tags comments
//// @Accept json
//// @Produce json
//// @Security Bearer
//// @Success 200 {array} domain.CommentResponse
//// @Failure 401 {object} domain.ErrorResponse
//// @Failure 500 {object} domain.ErrorResponse
//// @Router /api/v1/users/me/comments [get]
//func (h *Handler) GetUserComments(c *gin.Context) {
//	userId, err := getUserIdFromContext(c)
//	if err != nil {
//		h.logger.Errorw("Failed to get user Id from context", "error", err)
//		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
//			Error:   "Unauthorized",
//			Message: "You must be logged in to view your comments",
//		})
//		return
//	}
//
//	// Forward to communication service via proxy
//	c.Request.URL.Path = "/api/v1/comments/user/" + strconv.Itoa(int(userId))
//
//	proxy, err := h.service.Proxy.NewCommunicationServiceProxy()
//	if err != nil {
//		h.logger.Errorw("Failed to create communication service proxy", "error", err)
//		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
//			Error:   "Internal Server Error",
//			Message: "Failed to retrieve comments",
//		})
//		return
//	}
//
//	// Add user Id to request context
//	c.Request.Header.Set("X-User-Id", strconv.Itoa(int(userId)))
//
//	proxy.ServeHTTP(c.Writer, c.Request)
//}
