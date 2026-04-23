package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vinzryyy/iot-backend/models"
	"github.com/vinzryyy/iot-backend/repo"
)

var ErrForbidden = errors.New("no access to this location")

type DeviceService struct {
	devices   *repo.DeviceRepository
	users     *repo.UserRepository
	locations *repo.LocationRepository
}

func NewDeviceService(d *repo.DeviceRepository, u *repo.UserRepository, l *repo.LocationRepository) *DeviceService {
	return &DeviceService{devices: d, users: u, locations: l}
}

func (s *DeviceService) allowedLocations(ctx context.Context, userID, role string) ([]string, map[string]struct{}, error) {
	ids, err := s.users.LocationIDs(ctx, userID, role)
	if err != nil {
		return nil, nil, err
	}
	set := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return ids, set, nil
}

func (s *DeviceService) List(ctx context.Context, userID, role string) ([]models.Device, error) {
	ids, _, err := s.allowedLocations(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	return s.devices.List(ctx, ids)
}

func (s *DeviceService) Get(ctx context.Context, userID, role, id string) (*models.Device, error) {
	_, allowed, err := s.allowedLocations(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	d, err := s.devices.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if _, ok := allowed[d.LocationID]; !ok {
		return nil, ErrForbidden
	}
	return d, nil
}

func (s *DeviceService) Create(ctx context.Context, userID, role string, req models.CreateDeviceRequest) (*models.Device, error) {
	_, allowed, err := s.allowedLocations(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	if _, ok := allowed[req.LocationID]; !ok {
		return nil, ErrForbidden
	}

	exists, err := s.locations.Exists(ctx, req.LocationID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("location not found")
	}

	status := req.Status
	if status == "" {
		status = "offline"
	}

	d := &models.Device{
		ID:         uuid.NewString(),
		Name:       req.Name,
		LocationID: req.LocationID,
		Status:     status,
	}
	if err := s.devices.Create(ctx, d); err != nil {
		return nil, err
	}
	return s.devices.FindByID(ctx, d.ID)
}

func (s *DeviceService) Update(ctx context.Context, userID, role, id string, req models.UpdateDeviceRequest) (*models.Device, error) {
	_, allowed, err := s.allowedLocations(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	existing, err := s.devices.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if _, ok := allowed[existing.LocationID]; !ok {
		return nil, ErrForbidden
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.LocationID != "" && req.LocationID != existing.LocationID {
		if _, ok := allowed[req.LocationID]; !ok {
			return nil, ErrForbidden
		}
		exists, err := s.locations.Exists(ctx, req.LocationID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("location not found")
		}
		existing.LocationID = req.LocationID
	}
	if req.Status != "" {
		existing.Status = req.Status
	}
	if err := s.devices.Update(ctx, existing); err != nil {
		return nil, err
	}
	return s.devices.FindByID(ctx, existing.ID)
}

func (s *DeviceService) Delete(ctx context.Context, userID, role, id string) error {
	_, allowed, err := s.allowedLocations(ctx, userID, role)
	if err != nil {
		return err
	}
	existing, err := s.devices.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if _, ok := allowed[existing.LocationID]; !ok {
		return ErrForbidden
	}
	return s.devices.Delete(ctx, id)
}
