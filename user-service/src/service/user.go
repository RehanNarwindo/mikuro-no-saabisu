package service

import (
	"errors"
	"fmt"
	"user-service/src/dto"
	"user-service/src/repository"
	"github.com/golang-jwt/jwt/v5"
)


func GetPublicMessage() string {
	return "User service jalan"
}

func GetProfile(claims jwt.MapClaims) (dto.UserResponse, error) {
	userID, ok := claims["sub"].(string)
	if !ok {
		return dto.UserResponse{}, errors.New("invalid token: user_id not found")
	}

	user, err := repository.GetUserById(userID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name: user.FirstName + " " + user.LastName,
		Role:      user.Role,
	}, nil
}


func getFieldValue(data interface{}, fieldName string) interface{} {
	if m, ok := data.(map[string]interface{}); ok {
		if val, exists := m[fieldName]; exists {
			return val
		}
	}
	return nil
}

func GetUserByID(claims jwt.MapClaims, targetUserID string) (dto.UserResponse, error) {
	tokenUserID, _ := claims["sub"].(string)
	tokenRole, _ := claims["role"].(string)

	if tokenRole != "admin" && tokenRole != "super_admin" && tokenUserID != targetUserID {
		return dto.UserResponse{}, errors.New("permission denied: you can only access your own data")
	}

	user, err := repository.GetUserById(targetUserID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name: user.FirstName + " " + user.LastName,
		Role:      user.Role,
	}, nil
}




func GetAllUsers(claims jwt.MapClaims, req dto.GetAllUsersRequest) (*dto.GetAllUsersResponse, error) {
	tokenRole, _ := claims["role"].(string)
	if tokenRole != "admin" && tokenRole != "super_admin" {
		return nil, errors.New("permission denied: admin access required")
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortDir == "" {
		req.SortDir = "DESC"
	}

	users, total, err := repository.GetAllUsersWithFilters(
		req.Search,
		req.Role,
		req.Limit,
		req.Offset,
		req.SortBy,
		req.SortDir,
	)
	if err != nil {
		return nil, err
	}

	totalPages := (total + req.Limit - 1) / req.Limit
	if totalPages < 0 {
		totalPages = 0
	}

	var result []dto.UserResponse
	for _, user := range users {
		result = append(result, dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.FirstName + " " + user.LastName,
			Role:  user.Role,
		})
	}

	return &dto.GetAllUsersResponse{
		Users:      result,
		Total:      total,
		Limit:      req.Limit,
		Offset:     req.Offset,
		TotalPages: totalPages,
	}, nil
}


func UpdateUser(claims jwt.MapClaims, targetUserID string, updateData dto.UpdateUserRequest) (dto.UserResponse, error) {
	tokenUserID, _ := claims["sub"].(string)
	tokenRole, _ := claims["role"].(string)

	if tokenRole != "admin" && tokenRole != "super_admin" && tokenUserID != targetUserID {
		return dto.UserResponse{}, errors.New("permission denied: you can only update your own data")
	}

	existingUser, err := repository.GetUserById(targetUserID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	updates := make(map[string]interface{})
	if updateData.Email != "" && updateData.Email != existingUser.Email {
		updates["email"] = updateData.Email
	}
	if updateData.FirstName != "" {
		updates["first_name"] = updateData.FirstName
	}
	if updateData.LastName != "" {
		updates["last_name"] = updateData.LastName
	}

	err = repository.UpdateUserByID(targetUserID, updates)
	if err != nil {
		return dto.UserResponse{}, err
	}

	updatedUser, err := repository.GetUserById(targetUserID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		Name: 	   updatedUser.FirstName + " " + updatedUser.LastName,
		Role:      updatedUser.Role,
	}, nil
}

func DeleteUser(claims jwt.MapClaims, targetUserID string) error {
	tokenUserID, _ := claims["sub"].(string)
	tokenRole, _ := claims["role"].(string)

	fmt.Printf("DeleteUser - Token UserID: %s, Role: %s, Target UserID: %s\n", 
		tokenUserID, tokenRole, targetUserID)

	if tokenRole != "admin" && tokenRole != "super_admin" && tokenUserID != targetUserID {
		return errors.New("permission denied: you can only delete your own account")
	}

	if tokenUserID == targetUserID && tokenRole == "admin" {
		fmt.Println("Warning: Admin is deleting their own account")
	}

	err := repository.DeleteUserByID(targetUserID)
	if err != nil {
		return err
	}

	fmt.Printf("User with ID %s has been deleted by %s (role: %s)\n", 
		targetUserID, tokenUserID, tokenRole)
	
	return nil
}