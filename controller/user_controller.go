package controller

import (
	"example/dtos"
	apierrors "example/errors"
	"example/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	var req dtos.GetAllUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		resp := apierrors.ErrorBadRequest
		resp.Message = "invalid query parameters"
		apierrors.CommonErrorResponse(ctx, &resp)
		return
	}
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")

	// Convert page and perPage to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		resp := apierrors.ErrorBadRequest
		resp.Message = "invalid page number"
		apierrors.CommonErrorResponse(ctx, &resp)
		return
	}
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil || perPageInt < 1 {
		resp := apierrors.ErrorBadRequest
		resp.Message = "invalid per_page number"
		apierrors.CommonErrorResponse(ctx, &resp)
		return
	}

	// Calculate limit and offset for pagination
	limit := perPageInt
	offset := (pageInt - 1) * perPageInt

	users, total, resp := c.service.GetAllUsers(limit, offset)
	if resp != nil {
		apierrors.CommonErrorResponse(ctx, resp)
		return
	}
	apierrors.CommonSuccessResponse(ctx, dtos.CommonPaginationData{
		Data:   users,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	// Implementation for getting a user by ID

	userID := ctx.Param("id")
	idInt, err := strconv.Atoi(userID)
	if err != nil {
		resp := apierrors.ErrorBadRequest
		resp.Message = "invalid user id"
		apierrors.CommonErrorResponse(ctx, &resp)
		return
	}

	user, resp := c.service.GetUserByID(idInt)
	if resp != nil {
		apierrors.CommonErrorResponse(ctx, resp)
		return
	}

	apierrors.CommonSuccessResponse(ctx, user)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	// Implementation for creating a new user
	var req dtos.NewUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := apierrors.ErrorBadRequest
		resp.Message = "invalid request body"
		apierrors.CommonErrorResponse(ctx,& resp)
		return
	}

	createdUser, resp := c.service.CreateUser(
		dtos.UserDTO{
			Username: req.Username,
			Email:    req.Email,
			Role:     "user", // Default role for new users
		}) // convert request DTO to service DTO
	if resp != nil {
		apierrors.CommonErrorResponse(ctx, resp)
		return
	}
	apierrors.CommonSuccessResponse(ctx, createdUser)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	// Implementation for updating an existing user
	userID := ctx.Param("id")
	idInt, err := strconv.Atoi(userID)
	if err != nil {
		resp := apierrors.ErrorBadRequest
		resp.Message = "invalid user id"
		apierrors.CommonErrorResponse(ctx, &resp)
		return
	}

	var req dtos.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := apierrors.ErrorBadRequest
		resp.Message = "invalid request body"
		apierrors.CommonErrorResponse(ctx, &resp)
		return
	}

	updatedUser, resp := c.service.UpdateUser(idInt, dtos.UserDTO{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role, // Allow role update if provided
	})
	if resp != nil {
		apierrors.CommonErrorResponse(ctx, resp)
		return
	}
	apierrors.CommonSuccessResponse(ctx, updatedUser)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	// Implementation for deleting a user
	userID := ctx.Param("id")
	idInt, err := strconv.Atoi(userID)
	if err != nil {
		resp := apierrors.ErrorBadRequest
		resp.Message = "invalid user id"
		apierrors.CommonErrorResponse(ctx, &resp)
		return
	}

	if resp := c.service.DeleteUser(idInt); resp != nil {
		apierrors.CommonErrorResponse(ctx, resp)
		return
	}
	apierrors.CommonSuccessResponse(ctx, gin.H{"message": "user deleted successfully"})
}
