package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Token struct {
	Plaintext string    `json:"token" validate:"required,len=26"`
	Hash      [32]byte  `json:"-"`
	UserID    int32     `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

func (s *Token) Validate() error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}

	return nil
}
