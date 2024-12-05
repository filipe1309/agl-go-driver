package users

import (
	"time"

	domain "github.com/filipe1309/agl-go-driver/internal/users"
)

func (h *handler) authenticate(login, password string) (*domain.User, error) {
	return h.factory.Authenticate(login, password)
}

func (h *handler) updateLastLogin(user *domain.User) (int64, error) {
	user.LastLogin = time.Now()
	return h.repo.UpdateDB(user.ID, user)
}

func Authenticate(login, password string) (u *domain.User, err error) {
	u, err = gh.authenticate(login, password)
	if err != nil {
		return
	}

	_, err = gh.updateLastLogin(u)
	return
}
