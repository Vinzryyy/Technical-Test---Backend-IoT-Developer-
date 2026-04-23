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

// Register is public self-signup. It forces role=user and grants no
// locations — an admin must provision access afterwards.
func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	return s.createUser(ctx, req.Name, req.Email, req.Password, "user", nil)
}

// RegisterStaff is admin-only: creates a user with an explicit role and
// location access list.
func (s *AuthService) RegisterStaff(ctx context.Context, req models.StaffRegisterRequest) (*models.User, error) {
	role := req.Role
	if role == "" {
		role = "user"
	}
	return s.createUser(ctx, req.Name, req.Email, req.Password, role, req.LocationIDs)
}

func (s *AuthService) createUser(ctx context.Context, name, email, password, role string, locationIDs []string) (*models.User, error) {
	if existing, err := s.users.FindByEmail(ctx, email); err == nil && existing != nil {
		return nil, ErrEmailExists
	} else if err != nil && !errors.Is(err, repo.ErrNotFound) {
		return nil, err
	}

	for _, lid := range locationIDs {
		ok, err := s.locs.Exists(ctx, lid)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("unknown location id: " + lid)
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		ID:       uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: string(hash),
		Role:     role,
	}
	if err := s.users.Create(ctx, u, locationIDs); err != nil {
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
