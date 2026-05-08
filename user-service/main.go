package main

import (
	"context"
	"log"
	"time"

	"user-service/src/config/database"
	"user-service/src/config/jwt"
	handler "user-service/src/handler/user"
	"user-service/src/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := initDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	if err := runServer(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func runServer() error {
	defer func() {
		if err := database.CloseDatabase(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	cfg := jwt.LoadConfigJWT()

	authMiddleware := middleware.AuthMiddleware(middleware.Options{
		JWTSecret: cfg.JWTSecret,
	})

	r := gin.Default()

	r.GET("/public", handler.PublicHandler)
	r.GET("/users/all", authMiddleware, handler.GetAllUserHandler)
	r.GET("/users/profile", authMiddleware, handler.GetUserProfileHandler)
	r.GET("/users/profile/:id", authMiddleware, handler.GetUserByIdHandler)
	r.PUT("/users/:id", authMiddleware, handler.UpdateUserHandler)
	r.DELETE("/users/:id", authMiddleware, handler.DeleteUserHandler)

	return r.Run(":3001")
}

func initDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return database.DatabaseConfigWithRetryContext(ctx, 3)
}
