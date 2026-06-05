package routes

import (
	"example/controller"
	"example/middleware"
	"example/services"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, userController *controller.UserController, userService services.UserService) {
	userGroup := router.Group("/users")
	userGroup.Use(middleware.AuthMiddleware(userService)) // resolves identity + role for every /users route
	{
		userGroup.GET("/", middleware.RequireAdmin(), userController.GetAllUsers)
		userGroup.GET("/:id", middleware.RequireSelfOrAdmin(), userController.GetUserByID)
		userGroup.POST("/", middleware.RequireAdmin(), userController.CreateUser)
		userGroup.PUT("/:id", middleware.RequireSelfOrAdmin(), userController.UpdateUser)
		userGroup.DELETE("/:id", middleware.RequireAdmin(), userController.DeleteUser)
	}
}
