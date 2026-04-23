package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vinzryyy/iot-backend/models"
)

type DeviceRepository struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// List returns all devices whose location is in allowedLocationIDs.
// An empty slice means the caller has no access -> return empty list.
func (r *DeviceRepository) List(ctx context.Context, allowedLocationIDs []string) ([]models.Device, error) {
	if len(allowedLocationIDs) == 0 {
		return []models.Device{}, nil
	}

	placeholders := make([]string, len(allowedLocationIDs))
	args := make([]interface{}, len(allowedLocationIDs))
	for i, id := range allowedLocationIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	q := fmt.Sprintf(`
		SELECT d.id, d.name, d.location_id, l.name, d.status, d.updated_at, d.created_at
		FROM devices d
		JOIN locations l ON l.id = d.location_id
		WHERE d.location_id IN (%s)
		ORDER BY d.updated_at DESC`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Device, 0)
	for rows.Next() {
		var d models.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.LocationID, &d.Location,
			&d.Status, &d.UpdatedAt, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *DeviceRepository) FindByID(ctx context.Context, id string) (*models.Device, error) {
	const q = `
		SELECT d.id, d.name, d.location_id, l.name, d.status, d.updated_at, d.created_at
		FROM devices d
		JOIN locations l ON l.id = d.location_id
		WHERE d.id = $1`
	var d models.Device
	err := r.db.QueryRow(ctx, q, id).Scan(
		&d.ID, &d.Name, &d.LocationID, &d.Location,
		&d.Status, &d.UpdatedAt, &d.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DeviceRepository) Create(ctx context.Context, d *models.Device) error {
	d.CreatedAt = time.Now()
	d.UpdatedAt = d.CreatedAt
	const q = `
		INSERT INTO devices (id, name, location_id, status, updated_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, q,
		d.ID, d.Name, d.LocationID, d.Status, d.UpdatedAt, d.CreatedAt)
	return err
}

func (r *DeviceRepository) Update(ctx context.Context, d *models.Device) error {
	d.UpdatedAt = time.Now()
	const q = `
		UPDATE devices
		SET name = $1, location_id = $2, status = $3, updated_at = $4
		WHERE id = $5`
	ct, err := r.db.Exec(ctx, q, d.Name, d.LocationID, d.Status, d.UpdatedAt, d.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *DeviceRepository) Delete(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM devices WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
