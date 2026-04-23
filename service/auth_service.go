package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vinzryyy/iot-backend/models"
	"github.com/vinzryyy/iot-backend/repo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already registered")
)

type AuthService struct {
	users *repo.UserRepository
	locs  *repo.LocationRepository
	jwt   *JWTService
}

func NewAuthService(u *repo.UserRepository, l *repo.LocationRepository, j *JWTService) *AuthService {
	return &AuthService{users: u, locs: l, jwt: j}
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, expiresIn, err := s.jwt.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	locations, err := s.users.Locations(ctx, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User: models.UserProfile{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			Locations: locations,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	if existing, err := s.users.FindByEmail(ctx, req.Email); err == nil && existing != nil {
		return nil, ErrEmailExists
	} else if err != nil && !errors.Is(err, repo.ErrNotFound) {
		return nil, err
	}

	for _, lid := range req.LocationIDs {
		ok, err := s.locs.Exists(ctx, lid)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("unknown location id: " + lid)
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	u := &models.User{
		ID:       uuid.NewString(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
		Role:     role,
	}
	if err := s.users.Create(ctx, u, req.LocationIDs); err != nil {
		return nil, err
	}
	u.Password = ""
	return u, nil
}

func (s *AuthService) Profile(ctx context.Context, userID string) (*models.UserProfile, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	locations, err := s.users.Locations(ctx, u.ID, u.Role)
	if err != nil {
		return nil, err
	}
	return &models.UserProfile{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		Locations: locations,
		CreatedAt: u.CreatedAt,
	}, nil
}
