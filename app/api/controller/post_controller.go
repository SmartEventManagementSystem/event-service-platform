package controller

import (
    "context"
    "net/http"

    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/request"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
)

type PostController struct {
    res      runtime.Resource
    managers *manager.Managers
}

func NewPostController(managers *manager.Managers, res runtime.Resource) *PostController {
    return &PostController{res: res, managers: managers}
}

func (c *PostController) CreatePost(ec echo.Context) error {
    var req request.CreatePostRequest
    if err := ec.Bind(&req); err != nil {
        return ec.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
    }

    // extract current user id from context
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

    post, err := c.managers.PostManager.CreatePost(context.Background(), &req, userID)
    if err != nil {
        return ec.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create post"})
    }
    return ec.JSON(http.StatusCreated, post)
}

func (c *PostController) ListPosts(ec echo.Context) error {
    // list posts for a user; expect username param or user_id query
    userIDStr := ec.QueryParam("user_id")
    if userIDStr == "" {
        return ec.JSON(http.StatusBadRequest, map[string]string{"error": "user_id required"})
    }
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return ec.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
    }

    page := 1
    pageSize := 20
    posts, err := c.managers.PostManager.ListPostsByUser(context.Background(), userID, page, pageSize)
    if err != nil {
        return ec.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list posts"})
    }
    return ec.JSON(http.StatusOK, posts)
}

func (c *PostController) GetPost(ec echo.Context) error {
    idStr := ec.Param("id")
    if idStr == "" {
        return ec.JSON(http.StatusBadRequest, map[string]string{"error": "id required"})
    }
    id, err := uuid.Parse(idStr)
    if err != nil {
        return ec.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
    }
    post, err := c.managers.PostManager.GetPost(context.Background(), id)
    if err != nil {
        return ec.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get post"})
    }
    return ec.JSON(http.StatusOK, post)
}

func (c *PostController) UpdatePost(ec echo.Context) error {
    idStr := ec.Param("id")
    if idStr == "" {
        return ec.JSON(http.StatusBadRequest, map[string]string{"error": "id required"})
    }
    id, err := uuid.Parse(idStr)
    if err != nil {
        return ec.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
    }
    var req request.UpdatePostRequest
    if err := ec.Bind(&req); err != nil {
        return ec.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
    }
    post, err := c.managers.PostManager.UpdatePost(context.Background(), id, &req)
    if err != nil {
        return ec.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update post"})
    }
    return ec.JSON(http.StatusOK, post)
}

func (c *PostController) DeletePost(ec echo.Context) error {
    idStr := ec.Param("id")
    if idStr == "" {
        return ec.JSON(http.StatusBadRequest, map[string]string{"error": "id required"})
    }
    id, err := uuid.Parse(idStr)
    if err != nil {
        return ec.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
    }
    if err := c.managers.PostManager.DeletePost(context.Background(), id); err != nil {
        return ec.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete post"})
    }
    return ec.NoContent(http.StatusNoContent)
}
