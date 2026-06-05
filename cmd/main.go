package main

import (
	"example/configs"
	"example/controller"
	"example/repositories"
	"example/routes"
	"example/services"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load(".env")
	appName := "example-service"
	db := configs.SetupDatabase(&appName)

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "dev" || appEnv == "local" || appEnv == "sit" {
		mode := "up"
		steps := 1
		if len(os.Args) >= 2 {
			mode = os.Args[1]
		}
		if len(os.Args) >= 3 {
			n, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("Invalid step number: %v", err)
			}
			steps = n
		}
		configs.RunMigrations(db, mode, steps)
	} else {
		log.Println("Skipping migrations: APP_ENV is not 'dev', 'local' or 'sit'")
	}

	userRepo := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	userController := controller.NewUserController(userService)

	r := gin.Default()
	api := r.Group("/api/v1")

	routes.RegisterUserRoutes(api, userController, userService)

	r.Run(":8080")

}
