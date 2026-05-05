package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"user-service/src/config/database"
	"user-service/src/model"
)

func GetUserById(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	query := `SELECT id, email, first_name, last_name, role FROM users WHERE id = $1`

	err := database.DB.QueryRowContext(ctx, query, id).Scan(
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

func GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	query := `SELECT id, email, first_name, last_name, role FROM users WHERE email = $1`

	err := database.DB.QueryRowContext(ctx, query, email).Scan(
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

func GetAllUserHandlersWithFilters(ctx context.Context, search, role string, limit, offset int, sortBy, sortDir string) ([]model.User, int, error) {
	sortBy, sortDir = normalizeSortParams(sortBy, sortDir)
	limit, offset = normalizePaginationParams(limit, offset)

	baseQuery, args := buildBaseQuery(search, role)

	total, err := getTotalCount(ctx, baseQuery, args)
	if err != nil {
		return nil, 0, err
	}

	users, err := fetchUsers(ctx, baseQuery, args, sortBy, sortDir, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func normalizeSortParams(sortBy, sortDir string) (string, string) {
	allowedSortBy := map[string]bool{
		"id":         true,
		"email":      true,
		"first_name": true,
		"last_name":  true,
		"role":       true,
		"created_at": true,
	}

	if !allowedSortBy[sortBy] {
		sortBy = "created_at"
	}

	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "DESC"
	}

	return sortBy, sortDir
}

func normalizePaginationParams(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func buildBaseQuery(search, role string) (string, []any) {
	baseQuery := "FROM users WHERE 1=1"
	args := []any{}
	argCounter := 1

	if search != "" {
		baseQuery += fmt.Sprintf(` AND (first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d)`,
			argCounter, argCounter, argCounter)
		args = append(args, "%"+search+"%")
		argCounter++
	}

	if role != "" {
		baseQuery += fmt.Sprintf(` AND role = $%d`, argCounter)
		args = append(args, role)
	}

	return baseQuery, args
}

func getTotalCount(ctx context.Context, baseQuery string, args []any) (int, error) {
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := database.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	return total, err
}

func fetchUsers(ctx context.Context, baseQuery string, baseArgs []any, sortBy, sortDir string, limit, offset int) ([]model.User, error) {
	allowedSortBy := map[string]bool{
		"created_at": true,
		"email":      true,
		"first_name": true,
		"last_name":  true,
		"id":         true,
		"role":       true,
	}

	if !allowedSortBy[sortBy] {
		sortBy = "created_at"
	}

	sortDir = strings.ToUpper(sortDir)
	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "DESC"
	}

	// #nosec G201 - SQL injection prevented by whitelist validation above
	query := fmt.Sprintf(`
        SELECT id, email, first_name, last_name, role 
        %s 
        ORDER BY %s %s 
        LIMIT $%d OFFSET $%d
    `, baseQuery, sortBy, sortDir, len(baseArgs)+1, len(baseArgs)+2)

	args := append(baseArgs, limit, offset)

	rows, err := database.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Warning: failed to close rows: %v\n", err)
		}
	}()

	return scanUsers(rows)
}

func scanUsers(rows *sql.Rows) ([]model.User, error) {
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

func GetAllUserHandlers(ctx context.Context) ([]model.User, error) {
	rows, err := database.DB.QueryContext(ctx, `SELECT id, email, first_name, last_name, role FROM users`)
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

func UpdateUserByID(ctx context.Context, id string, updates map[string]any) error {
	query := `UPDATE users SET 
		email = COALESCE($1, email),
		first_name = COALESCE($2, first_name),
		last_name = COALESCE($3, last_name),
		updated_at = NOW()
	WHERE id = $4`

	_, err := database.DB.ExecContext(ctx, query,
		updates["email"],
		updates["first_name"],
		updates["last_name"],
		id,
	)

	return err
}

func DeleteUserByID(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := database.DB.ExecContext(ctx, query, id)
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
