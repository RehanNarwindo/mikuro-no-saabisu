package service

import (
	"errors"
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
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}, nil
}

func GetUserByID(claims jwt.MapClaims, targetUserID string) (dto.UserResponse, error) {
	tokenUserID, _ := claims["sub"].(string)
	tokenRole, _ := claims["role"].(string)

	if tokenRole != "admin" && tokenUserID != targetUserID {
		return dto.UserResponse{}, errors.New("permission denied: you can only access your own data")
	}

	user, err := repository.GetUserById(targetUserID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}, nil
}

func GetAllUsers(claims jwt.MapClaims, req dto.GetAllUserHandlersRequest) (*dto.GetAllUsersResponse, error) {
	tokenRole, _ := claims["role"].(string)
	if tokenRole != "admin" {
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

	users, total, err := repository.GetAllUserHandlersWithFilters(
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

	totalPages := 0
	if req.Limit > 0 {
		totalPages = (total + req.Limit - 1) / req.Limit
	}

	var result []dto.UserResponse
	for _, user := range users {
		result = append(result, dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
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

	if tokenRole != "admin" && tokenUserID != targetUserID {
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

	if len(updates) == 0 {
		return dto.UserResponse{
			ID:        existingUser.ID,
			Email:     existingUser.Email,
			FirstName: existingUser.FirstName,
			LastName:  existingUser.LastName,
			Role:      existingUser.Role,
		}, nil
	}

	if err := repository.UpdateUserByID(targetUserID, updates); err != nil {
		return dto.UserResponse{}, err
	}

	updatedUser, err := repository.GetUserById(targetUserID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Role:      updatedUser.Role,
	}, nil
}

func DeleteUser(claims jwt.MapClaims, targetUserID string) error {
	tokenUserID, _ := claims["sub"].(string)
	tokenRole, _ := claims["role"].(string)

	if tokenRole != "admin" && tokenUserID != targetUserID {
		return errors.New("permission denied: you can only delete your own account")
	}

	return repository.DeleteUserByID(targetUserID)
}
