package handler

import (
	"net/http"

	"user-service/src/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func PublicHandler(c *gin.Context) {
	message := service.GetPublicMessage()

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func ProfileHandler(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Invalid claims format",
		})
		return
	}

	user, err := service.GetProfile(claims)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    user,
	})
}

func GetAllUser(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	_, ok := userClaims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Invalid claims format",
		})
		return
	}

	users, err := service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    users,
	})
}