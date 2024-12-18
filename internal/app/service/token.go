package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"log/slog"
	"time"

	"activator/internal/app/models"
	"activator/internal/config"
)

type TokenRepository interface {
	Save(user models.Token) error
	DeleteAll(id int32) error
}

type TokenService struct {
	cfg    config.Config
	logger *slog.Logger
	repo   TokenRepository
}

func NewTokenService(cfg config.Config, logger *slog.Logger, repo TokenRepository) *TokenService {
	return &TokenService{
		cfg:    cfg,
		logger: logger,
		repo:   repo,
	}
}

func (s *TokenService) Create(user models.User) (models.Token, error) {
	// if err := params.Validate(); err != nil {
	// 	return models.User{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "service create")
	// }

	newToken, err := generateToken(user.ID, 3*24*time.Hour)
	if err != nil {
		return models.Token{}, err
	}

	// token := fmt.Sprintf("%b", newToken.Hash)
	// s.logger.Debug(token)

	err = s.repo.Save(newToken)
	if err != nil {
		return models.Token{}, err
	}

	return newToken, nil
}

func (s *TokenService) DeleteAll(id int32) error {
	if err := s.repo.DeleteAll(id); err != nil {
		return err
	}

	return nil
}

func generateToken(userId int32, ttl time.Duration) (models.Token, error) {
	token := models.Token{
		UserID: userId,
		Expiry: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return models.Token{}, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	token.Hash = sha256.Sum256([]byte(token.Plaintext))

	return token, nil
}
