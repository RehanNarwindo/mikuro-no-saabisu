package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"user-service/src/config/database"
	"user-service/src/model"
)

var allowedSortBy = map[string]bool{
	"created_at": true, "email": true, "first_name": true,
}

var allowedSortDir = map[string]bool{
	"ASC": true, "DESC": true,
}

func GetUserById(id string) (*model.User, error) {
	var user model.User

	query := `SELECT id, email, first_name, last_name, role FROM users WHERE id = $1`

	err := database.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User

	query := `SELECT id, email, first_name, last_name, role FROM users WHERE email = $1`

	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUserHandlersWithFilters(search, role string, limit, offset int, sortBy, sortDir string) ([]model.User, int, error) {
	baseQuery := `FROM users WHERE 1=1`
	args := []interface{}{}
	argCounter := 1

	if !allowedSortBy[sortBy] {
		sortBy = "created_at"
	}
	if !allowedSortDir[sortDir] {
		sortDir = "DESC"
	}

	if search != "" {
		baseQuery += fmt.Sprintf(` AND (first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d)`,
			argCounter, argCounter, argCounter)
		args = append(args, "%"+search+"%")
		argCounter++
	}

	if role != "" {
		baseQuery += fmt.Sprintf(` AND role = $%d`, argCounter)
		args = append(args, role)
		argCounter++
	}

	countQuery := `SELECT COUNT(*) ` + baseQuery
	var total int
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortDir == "" {
		sortDir = "DESC"
	}
	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "DESC"
	}

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT id, email, first_name, last_name, role 
		%s 
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d
	`, baseQuery, sortBy, sortDir, argCounter, argCounter+1)

	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Warning: failed to close rows: %v\n", err)
		}
	}()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.Role,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}

func GetAllUserHandlers() ([]model.User, error) {
	rows, err := database.DB.Query(`SELECT id, email, first_name, last_name, role FROM users`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Warning: failed to close rows: %v\n", err)
		}
	}()

	
	var users []model.User

	for rows.Next() {
		var u model.User

		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.Role,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func UpdateUserByID(id string, updates map[string]interface{}) error {
	query := `UPDATE users SET 
		email = COALESCE($1, email),
		first_name = COALESCE($2, first_name),
		last_name = COALESCE($3, last_name),
		updated_at = NOW()
	WHERE id = $4`

	_, err := database.DB.Exec(query,
		updates["email"],
		updates["first_name"],
		updates["last_name"],
		id,
	)

	return err
}

func DeleteUserByID(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
