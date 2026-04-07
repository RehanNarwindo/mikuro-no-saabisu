package main
import (
	"user-service/src/config/database"
	"user-service/src/config/jwt"
	"user-service/src/handler"
	"user-service/src/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.DatabaseConfig()

	cfg := jwt.LoadConfigJWT()

	authMiddleware := middleware.AuthMiddleware(middleware.Options{
		JWTSecret: cfg.JWTSecret,
	})

	r := gin.Default()

	r.GET("/public", handler.PublicHandler)
	r.GET("/profile", authMiddleware, handler.ProfileHandler)

	r.Run(":3001")
}