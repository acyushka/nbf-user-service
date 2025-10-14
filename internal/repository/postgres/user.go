package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"nbf-user/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (id, name, surname, contacts, description)
        VALUES ($1, $2, $3, $4, $5)
    `

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Surname,
		pq.Array(user.Contacts),
		user.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	var contacts []string

	query := `SELECT id, name, surname, contacts, description
	          FROM users WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		pq.Array(&contacts),
		&user.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User not found: %w", err)
		}
		return nil, fmt.Errorf("Failed to get user: %w", err)
	}

	user.Contacts = contacts
	return &user, nil
}

func (r *UserRepo) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.User, error) {
	if len(ids) == 0 {
		return []*models.User{}, nil
	}

	type userRow struct {
		ID          uuid.UUID `db:"id"`
		Name        string    `db:"name"`
		Surname     string    `db:"surname"`
		Contacts    []byte    `db:"contacts"`
		Description string    `db:"description"`
	}

	query, args, err := sqlx.In(`
		SELECT id, name, surname, contacts, description
		FROM users
		WHERE id IN (?)
		ORDER BY name, surname
	`, ids)
	if err != nil {
		return nil, fmt.Errorf("Failed to build query: %w", err)
	}
	query = r.db.Rebind(query)
	var rows []userRow
	if err = r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("Failed to get users by ids: %w", err)
	}

	users := make([]*models.User, 0, len(rows))
	for _, row := range rows {
		contacts, err := parsePostgresArray(row.Contacts)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse contacts for user %s: %w", row.ID, err)
		}

		users = append(users, &models.User{
			ID:          row.ID,
			Name:        row.Name,
			Surname:     row.Surname,
			Contacts:    contacts,
			Description: row.Description,
		})
	}

	return users, nil
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, surname = $2, contacts = $3, description = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Surname,
		pq.Array(user.Contacts),
		user.Description,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("Failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("User not found")
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("User not found")
	}

	return nil
}

func parsePostgresArray(arrayData []byte) ([]string, error) {
	if len(arrayData) == 0 {
		return []string{}, nil
	}

	var result []string
	if err := pq.Array(&result).Scan(arrayData); err != nil {
		return nil, fmt.Errorf("Failed to parse postgres array: %w", err)
	}

	return result, nil
}
