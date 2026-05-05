package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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
		"status":  true,
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

	if validationErr := validate.Var(req.ID, "required,uuid"); validationErr != nil {
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

	req := parseUserRequest(c)
	if req == nil {
		return
	}

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

func parseUserRequest(c *gin.Context) *dto.GetAllUserHandlersRequest {
	req := &dto.GetAllUserHandlersRequest{
		Search:  c.Query("search"),
		Role:    c.Query("role"),
		Limit:   10,
		Offset:  0,
		SortBy:  "created_at",
		SortDir: "DESC",
	}

	if limit := c.Query("limit"); limit != "" {
		val, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid limit"})
			return nil
		}
		req.Limit = val
	}

	if offset := c.Query("offset"); offset != "" {
		val, err := strconv.Atoi(offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid offset"})
			return nil
		}
		req.Offset = val
	}

	allowedSortBy := map[string]bool{
		"created_at": true, "email": true, "first_name": true,
	}
	if sortBy := c.Query("sort_by"); sortBy != "" {
		if !allowedSortBy[sortBy] {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid sort_by"})
			return nil
		}
		req.SortBy = sortBy
	}

	if sortDir := c.Query("sort_dir"); sortDir != "" {
		sortDirUpper := strings.ToUpper(sortDir)
		if sortDirUpper != "ASC" && sortDirUpper != "DESC" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid sort_dir"})
			return nil
		}
		req.SortDir = sortDirUpper
	}

	return req
}

func UpdateUserHandler(c *gin.Context) {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	targetUserID := c.Param("id")

	var updateData dto.UpdateUserRequest
	if err = c.ShouldBindJSON(&updateData); err != nil {
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
		return nil, errors.New("unauthorized")
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}

	return claims, nil
}
