package services

import (
	"time"

	"github.com/filipe1309/agl-go-driver/factories"
	domain "github.com/filipe1309/agl-go-driver/internal/users"
	"github.com/filipe1309/agl-go-driver/repositories"
)

func NewAuthService(repo repositories.UserWriteRepository, factory *factories.UserFactory) *AuthService {
	return &AuthService{repo, factory}
}

type AuthService struct {
	repo    repositories.UserWriteRepository
	factory *factories.UserFactory
}

func (svc *AuthService) authenticate(login, password string) (*domain.User, error) {
	return svc.factory.Authenticate(login, password)
}

func (svc *AuthService) updateLastLogin(user *domain.User) (int64, error) {
	user.LastLogin = time.Now()
	return svc.repo.UpdateDB(user.ID, user)
}

func (svc *AuthService) Authenticate(login, password string) (u *domain.User, err error) {
	u, err = svc.authenticate(login, password)
	if err != nil {
		return
	}

	_, err = svc.updateLastLogin(u)
	return
}
