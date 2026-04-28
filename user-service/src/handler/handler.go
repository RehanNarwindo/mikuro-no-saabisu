package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"user-service/src/dto"
	"user-service/src/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

var validate = validator.New()


func PublicHandler(c *gin.Context) {
	message := service.GetPublicMessage()
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}


func GetUserProfileHandler(c *gin.Context) {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	user, err := service.GetProfile(claims)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"message": "Success",
		"data":    user,
	})
}

func GetUserByIdHandler(c *gin.Context) {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var req dto.GetUserByIdRequest
	req.ID = c.Param("id")
	
	if err := validate.Var(req.ID, "required,uuid"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID format"})
		return
	}

	user, err := service.GetUserByID(claims, req.ID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User retrieved successfully",
		"data":    user,
	})
}



func GetAllUserHandler(c *gin.Context) {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var req dto.GetAllUserHandlersRequest
	
	req.Search = c.Query("search")
	req.Role = c.Query("role")
	
	if limit := c.Query("limit"); limit != "" {
		req.Limit, _ = strconv.Atoi(limit)
	}
	if offset := c.Query("offset"); offset != "" {
		req.Offset, _ = strconv.Atoi(offset)
	}
	
	req.SortBy = c.Query("sort_by")
	req.SortDir = c.Query("sort_dir")

	response, err := service.GetAllUsers(claims, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Success",
		"data":    response,
	})
}

func UpdateUserHandler(c *gin.Context) {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is required"})
		return
	}

	var updateData dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	updatedUser, err := service.UpdateUser(claims, targetUserID, updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
		"data":    updatedUser,
	})
}

func DeleteUserHandler(c *gin.Context) {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is required"})
		return
	}

	err = service.DeleteUser(claims, targetUserID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

func getClaimsFromContext(c *gin.Context) (jwt.MapClaims, error) {
	userClaims, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims format")
	}

	return claims, nil
}