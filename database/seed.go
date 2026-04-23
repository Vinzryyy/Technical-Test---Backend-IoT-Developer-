package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	adminID      = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	jakartaUID   = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	surabayaUID  = "cccccccc-cccc-cccc-cccc-cccccccccccc"
	bandungUID   = "dddddddd-dddd-dddd-dddd-dddddddddddd"
	supervisorID = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"

	jakartaLocID  int64 = 1
	surabayaLocID int64 = 2
	bandungLocID  int64 = 3
)

type seedUser struct {
	id       string
	name     string
	email    string
	password string
	role     string
	// locations the user can access (ignored for admin)
	locations []int64
}

// SeedUsers inserts the default demo users on first run. Idempotent.
func SeedUsers(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	users := []seedUser{
		{
			id:       adminID,
			name:     "Super Admin",
			email:    "admin@example.com",
			password: "admin123",
			role:     "admin",
		},
		{
			id:        jakartaUID,
			name:      "Jakarta Staff",
			email:     "jakarta@example.com",
			password:  "user123",
			role:      "user",
			locations: []int64{jakartaLocID},
		},
		{
			id:        surabayaUID,
			name:      "Surabaya Staff",
			email:     "surabaya@example.com",
			password:  "user123",
			role:      "user",
			locations: []int64{surabayaLocID},
		},
		{
			id:        bandungUID,
			name:      "Bandung Staff",
			email:     "bandung@example.com",
			password:  "user123",
			role:      "user",
			locations: []int64{bandungLocID},
		},
		{
			id:        supervisorID,
			name:      "Regional Supervisor",
			email:     "supervisor@example.com",
			password:  "super123",
			role:      "user",
			locations: []int64{jakartaLocID, surabayaLocID},
		},
	}

	for _, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = pool.Exec(ctx, `
			INSERT INTO users (id, name, email, password, role)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (email) DO NOTHING`,
			u.id, u.name, u.email, string(hash), u.role)
		if err != nil {
			return err
		}

		for _, loc := range u.locations {
			_, err = pool.Exec(ctx, `
				INSERT INTO user_locations (user_id, location_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING`,
				u.id, loc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
