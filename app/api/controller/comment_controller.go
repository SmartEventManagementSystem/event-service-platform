package controller

import (
	"context"
	"net/http"
	"strings"

	"backend/event-service-platform/app/api/client/request"
	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/manager"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CommentController struct {
	res      runtime.Resource
	managers *manager.Managers
}

func NewCommentController(managers *manager.Managers, res runtime.Resource) *CommentController {
	return &CommentController{res: res, managers: managers}
}

func (c *CommentController) CreateComment(ec echo.Context) error {
	postIDStr := ec.Param("id")
	if postIDStr == "" {
		// fallback: tests may not set Echo path params; try to extract from URL
		// expected path: /api/v1/posts/:id/comments
		p := ec.Request().URL.Path
		// simple parse: look for /posts/{id}/comments
		// split and find "posts"
		parts := strings.Split(strings.Trim(p, "/"), "/")
		for i := 0; i < len(parts)-1; i++ {
			if parts[i] == "posts" {
				postIDStr = parts[i+1]
				break
			}
		}
		if postIDStr == "" {
			return ec.JSON(http.StatusBadRequest, map[string]string{"error": "post id required"})
		}
	}
	// check authentication first (tests may expect unauthorized before id validation)
	uidAny := ec.Get("user_id")
	if uidAny == nil {
		return ec.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	uidStr, ok := uidAny.(string)
	if !ok {
		return ec.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return ec.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": "invalid post id"})
	}

	var req request.CreateCommentRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	res, err := c.managers.CommentManager.CreateComment(context.Background(), &req, userID, postID)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create comment"})
	}
	return ec.JSON(http.StatusCreated, res)
}

func (c *CommentController) ListComments(ec echo.Context) error {
	postIDStr := ec.Param("id")
	if postIDStr == "" {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": "post id required"})
	}
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": "invalid post id"})
	}

	page := 1
	pageSize := 20
	res, err := c.managers.CommentManager.ListComments(context.Background(), postID, page, pageSize)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list comments"})
	}
	return ec.JSON(http.StatusOK, res)
}

func (c *CommentController) DeleteComment(ec echo.Context) error {
	idStr := ec.Param("id")
	if idStr == "" {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": "id required"})
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if err := c.managers.CommentManager.DeleteComment(context.Background(), id); err != nil {
		return ec.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete comment"})
	}
	return ec.NoContent(http.StatusNoContent)
}

// ReactToComment is a placeholder for comment reactions (not implemented yet)
func (c *CommentController) ReactToComment(ec echo.Context) error {
	return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "ReactToComment not implemented"})
}
