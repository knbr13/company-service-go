package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/knbr13/company-service-go/internal/repositories"
	"github.com/knbr13/company-service-go/internal/validator"
	"github.com/knbr13/company-service-go/pkg/util"
)

type UserService struct {
	Repos *repositories.Repositories
}

func (s *UserService) Register(ctx context.Context, user *repositories.User) error {
	v := validator.New()
	v.Check(len(user.Username) >= 3 && len(user.Username) <= 64, "username", "must be between 3 and 64 characters")
	v.Check(util.ValidMail(user.Email), "email", "must be a valid email address")
	v.Check(len(user.Password) >= 8, "password", "must at least contain 8 characters")

	if !v.Valid() {
		for k, v := range v.Errors {
			return validator.ValidationError(fmt.Sprintf("%s: %s", k, v))
		}
	}

	user.ID = uuid.New().String()

	return s.Repos.User.Insert(user)
}

func (s *UserService) Login(ctx context.Context, user *repositories.User) error {
	existingUser, err := s.Repos.User.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}

	ok, err := user.Matches([]byte(existingUser.Password))
	if err != nil {
		return err
	}
	if !ok {
		return repositories.ErrInvalidPassword
	}

	return nil
}
