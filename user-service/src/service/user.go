package service

import (
	"user-service/src/dto"
	"user-service/src/repository"
)


func GetPublicMessage() string {
	return "User service jalan"
}

func GetProfile(claims map[string]any) (dto.UserResponse, error) {
	email, _ := claims["email"].(string)

	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.FirstName + " " + user.LastName,
	}, nil
}

func GetAllUsers() ([]dto.UserResponse, error) {
	users, err := repository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var result []dto.UserResponse

	for _, u := range users {
		result = append(result, dto.UserResponse{
			ID:    u.ID,
			Email: u.Email,
			Name:  u.FirstName + " " + u.LastName,
		})
	}

	return result, nil
}