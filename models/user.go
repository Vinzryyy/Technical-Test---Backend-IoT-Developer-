package models

import "time"

type User struct {
	ID        string    `json:"id"         db:"id"`
	Name      string    `json:"name"       db:"name"`
	Email     string    `json:"email"      db:"email"`
	Password  string    `json:"-"          db:"password"`
	Role      string    `json:"role"       db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserProfile is what is returned by /me — user + accessible locations.
type UserProfile struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Locations []Location `json:"locations"`
	CreatedAt time.Time  `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	ExpiresIn   int64       `json:"expires_in"`
	User        UserProfile `json:"user"`
}

// RegisterRequest is used for public self-registration. Role is always
// forced to "user" and no location access is granted — an admin must
// provision locations afterwards.
type RegisterRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// StaffRegisterRequest is used by admins to create new staff with an
// explicit role and location access list.
type StaffRegisterRequest struct {
	Name        string  `json:"name"         validate:"required,min=2,max=100"`
	Email       string  `json:"email"        validate:"required,email"`
	Password    string  `json:"password"     validate:"required,min=6"`
	Role        string  `json:"role"         validate:"omitempty,oneof=admin user"`
	LocationIDs []int64 `json:"location_ids" validate:"omitempty,dive,gt=0"`
}
