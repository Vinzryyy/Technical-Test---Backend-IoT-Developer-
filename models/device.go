package models

import "time"

type Device struct {
	ID         int64     `json:"id"          db:"id"`
	Name       string    `json:"name"        db:"name"`
	LocationID int64     `json:"location_id" db:"location_id"`
	Location   string    `json:"location"    db:"location_name"`
	Status     string    `json:"status"      db:"status"`
	UpdatedAt  time.Time `json:"updated_at"  db:"updated_at"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
}

type CreateDeviceRequest struct {
	Name       string `json:"name"        validate:"required,min=1,max=150"`
	LocationID int64  `json:"location_id" validate:"required,gt=0"`
	Status     string `json:"status"      validate:"omitempty,oneof=online offline"`
}

type UpdateDeviceRequest struct {
	Name       string `json:"name"        validate:"omitempty,min=1,max=150"`
	LocationID int64  `json:"location_id" validate:"omitempty,gt=0"`
	Status     string `json:"status"      validate:"omitempty,oneof=online offline"`
}
