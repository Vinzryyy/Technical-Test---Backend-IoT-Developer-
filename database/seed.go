package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsers inserts default admin + a Jakarta-scoped staff user on first run.
func SeedUsers(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	adminHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userHash, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	const adminID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	const userID = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	const jakartaID = "11111111-1111-1111-1111-111111111111"

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, name, email, password, role)
		VALUES ($1, 'Super Admin', 'admin@example.com', $2, 'admin')
		ON CONFLICT (email) DO NOTHING`,
		adminID, string(adminHash))
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, name, email, password, role)
		VALUES ($1, 'Jakarta Staff', 'jakarta@example.com', $2, 'user')
		ON CONFLICT (email) DO NOTHING`,
		userID, string(userHash))
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO user_locations (user_id, location_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`,
		userID, jakartaID)
	return err
}
