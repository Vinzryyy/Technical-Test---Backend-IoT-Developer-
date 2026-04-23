package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vinzryyy/iot-backend/models"
)

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) List(ctx context.Context) ([]models.Location, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, created_at FROM locations ORDER BY name`)
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

func (r *LocationRepository) Exists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM locations WHERE id = $1)`, id).Scan(&exists)
	return exists, err
}
