package factories

import (
	"log"

	"github.com/filipe1309/agl-go-driver/internal/users"
	"github.com/filipe1309/agl-go-driver/repositories"
)

type scan interface {
	Scan(dest ...any) error
}

type UserFactory struct {
	repo repositories.UserReadRepository
}

func NewUserFactory(repo repositories.UserReadRepository) *UserFactory {
	return &UserFactory{repo}
}

func restore(row scan) (*users.User, error) {
	var user users.User
	err := row.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.Deleted)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (uf *UserFactory) RestoreOne(id int64) (*users.User, error) {
	row := uf.repo.ReadDB(id)
	return restore(row)
}

func (uf *UserFactory) RestoreAll() ([]users.User, error) {
	rows := uf.repo.ReadAllDB()
	users := make([]users.User, 0)
	for rows.Next() {
		user, err := restore(rows)
		if err != nil {
			log.Println(err)
			continue
		}

		users = append(users, *user)
	}

	return users, nil
}

func (uf *UserFactory) Authenticate(login, password string) (*users.User, error) {
	user := &users.User{Login: login}
	user.ChangePassword(password)

	row := uf.repo.Login(user.GetLogin(), user.GetPass())
	err := row.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.Deleted)
	if err != nil {
		return nil, err
	}

	return restore(row)
}
