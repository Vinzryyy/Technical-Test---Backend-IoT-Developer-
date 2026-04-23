package repo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vinzryyy/iot-backend/models"
)

var ErrNotFound = errors.New("resource not found")

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	const q = `
		SELECT id, name, email, password, role, created_at, updated_at
		FROM users WHERE email = $1`
	var u models.User
	err := r.db.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	const q = `
		SELECT id, name, email, password, role, created_at, updated_at
		FROM users WHERE id = $1`
	var u models.User
	err := r.db.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *models.User, locationIDs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt

	const insertUser = `
		INSERT INTO users (id, name, email, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := tx.Exec(ctx, insertUser,
		u.ID, u.Name, u.Email, u.Password, u.Role, u.CreatedAt, u.UpdatedAt,
	); err != nil {
		return err
	}

	for _, lid := range locationIDs {
		if _, err := tx.Exec(ctx,
			`INSERT INTO user_locations (user_id, location_id) VALUES ($1, $2)
			 ON CONFLICT DO NOTHING`,
			u.ID, lid); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// LocationIDs returns every location the user is allowed to access.
// Admins implicitly access every location, so this returns every location.
func (r *UserRepository) LocationIDs(ctx context.Context, userID, role string) ([]string, error) {
	var rows pgx.Rows
	var err error
	if role == "admin" {
		rows, err = r.db.Query(ctx, `SELECT id FROM locations`)
	} else {
		rows, err = r.db.Query(ctx,
			`SELECT location_id FROM user_locations WHERE user_id = $1`, userID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (r *UserRepository) Locations(ctx context.Context, userID, role string) ([]models.Location, error) {
	var rows pgx.Rows
	var err error
	if role == "admin" {
		rows, err = r.db.Query(ctx,
			`SELECT id, name, created_at FROM locations ORDER BY name`)
	} else {
		rows, err = r.db.Query(ctx, `
			SELECT l.id, l.name, l.created_at
			FROM locations l
			JOIN user_locations ul ON ul.location_id = l.id
			WHERE ul.user_id = $1
			ORDER BY l.name`, userID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Location, 0)
	for rows.Next() {
		var l models.Location
		if err := rows.Scan(&l.ID, &l.Name, &l.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}
