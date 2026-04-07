package repository

import (
	"user-service/src/config/database"
	"user-service/src/model"
)

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User

	query := `SELECT id, email, first_name, last_name FROM users WHERE email = $1`

	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]model.User, error) {
	rows, err := database.DB.Query(`SELECT id, email, first_name, last_name FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var u model.User

		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}