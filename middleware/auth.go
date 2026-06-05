package middleware

import (
	apierrors "example/errors"
	"net/http"
	"strconv"

	"example/services"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserID = "userID"
	ctxRole   = "role"
)

// AuthMiddleware resolves the caller from the X-User-Id header by loading the
// user from the database and stashing their id and role in the context.
// (Dev-grade identity; replace with real auth later without touching the guards.)
func AuthMiddleware(service services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("X-User-Id")
		if raw == "" {
			resp := apierrors.ErrorUnauthorized
			resp.Message = "X-User-Id header is required"
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}
		id, err := strconv.Atoi(raw)
		if err != nil {
			resp := apierrors.ErrorUnauthorized
			resp.Message = "X-User-Id must be an integer"
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}
		user, resp := service.GetUserByID(id)
		if resp != apierrors.SuccessResponse || user == nil {
			resp := apierrors.ErrorUnauthorized
			resp.Message = "unknown user"
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}
		c.Set(ctxUserID, id)
		c.Set(ctxRole, user.Role) // role from DB, never trusted from the client
		c.Next()
	}
}

// RequireAdmin allows only callers whose resolved role is "admin".
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(ctxRole) != "admin" {
			resp := apierrors.ErrorForbidden
			resp.Message = "admin role required"
			c.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}
		c.Next()
	}
}

// RequireSelfOrAdmin allows admins, or a user acting on their own :id record.
func RequireSelfOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(ctxRole) == "admin" {
			c.Next()
			return
		}
		paramID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			resp := apierrors.ErrorBadRequest
			resp.Message = "invalid user id"
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}
		if c.GetInt(ctxUserID) != paramID {
			resp := apierrors.ErrorForbidden
			resp.Message = "access denied"
			c.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}
		c.Next()
	}
}
