package service

import (
	"activator/internal"
	"activator/internal/app/models"
	"activator/internal/config"
	m "activator/internal/rest/models"
	"log/slog"
)

type UserRepository interface {
	Create(params m.CreateParams) (models.User, error)
	Delete(id int32) error
	FindByToken(tokenPlaintext string) (int32, error)
	Activate(id int32) error
	FindById(id int32) (models.User, error)
}

type UserService struct {
	cfg    config.Config
	logger *slog.Logger
	repo   UserRepository
}

func NewUserService(cfg config.Config, logger *slog.Logger, repo UserRepository) *UserService {
	return &UserService{
		cfg:    cfg,
		logger: logger,
		repo:   repo,
	}
}

func (s *UserService) Create(params m.CreateParams) (models.User, error) {
	if err := params.Validate(); err != nil {
		return models.User{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "service create")
	}

	newUser, err := s.repo.Create(params)
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

func (s *UserService) Delete(id int32) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Activate(tokenPlaintext string) (int32, error) {
	id, err := s.repo.FindByToken(tokenPlaintext)
	if err != nil {
		return 0, err
	}

	if err := s.repo.Activate(id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserService) Find(id int32) (models.User, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
