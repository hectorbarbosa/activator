package models

import (
	"github.com/go-playground/validator/v10"
)

type CreateParams struct {
	Email    string `json:"email" validate:"required"`
	Name     string `json:"user_name" validate:"required"`
	NickName string `json:"nick_name" validate:"required"`
}

func (s *CreateParams) Validate() error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}

	return nil
}

type UpdateParams struct {
	Email    string `json:"email" validate:"required"`
	Name     string `json:"user_name" validate:"required"`
	NickName string `json:"nick_name" validate:"required"`
}

func (s *UpdateParams) Validate() error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}

	return nil
}
