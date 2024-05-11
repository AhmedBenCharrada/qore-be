package person

import (
	"context"
	"fmt"
	"log/slog"
	"qore-be/internal/domain/dto"
)

// Repository the person repository.
//
//go:generate mockery --name=Repository --structname=PersonRepository  --case underscore --output=../mocks/ --filename=person_repository.go
type Repository interface {
	Add(context.Context, dto.PersonDTO) (dto.PersonDTO, error)
	GetByID(context.Context, int) (dto.PersonDTO, error)
}

// ServiceImpl implements the person service.
type ServiceImpl struct {
	db  Repository
	log *slog.Logger
}

// ServiceOption ..
type ServiceOption func(*ServiceImpl) error

// NewService creates a new person service.
func NewService(opts ...ServiceOption) (*ServiceImpl, error) {
	svc := &ServiceImpl{}

	for _, opt := range opts {
		err := opt(svc)
		if err != nil {
			return nil, err
		}
	}

	if svc.db == nil {
		return nil, fmt.Errorf("missing person DB")
	}

	if svc.log == nil {
		svc.log = slog.Default()
	}

	return svc, nil
}

// WithRepository returns a closure that initialize the person-service with the person-repository.
func WithRepository(db Repository) ServiceOption {
	return func(svc *ServiceImpl) error {
		svc.db = db
		return nil
	}
}

// Create saves a new person to database.
func (s *ServiceImpl) Create(ctx context.Context, d dto.PersonDTO) (dto.PersonDTO, error) {
	p, err := s.db.Add(ctx, d)
	if err != nil {
		s.log.Error("failed to save the person", "error", err.Error())
		return dto.PersonDTO{}, err
	}

	s.log.Info("person created with success", "person", p)
	return p, nil
}

// GetByID retrieves a person from the database by its ID.
func (s *ServiceImpl) GetByID(ctx context.Context, id int) (*dto.PersonDTO, error) {
	p, err := s.db.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get the person by id", "id", id, "error", err.Error())
		return nil, err
	}

	s.log.Info("person retrieved with success", "person", p)
	return &p, nil
}
