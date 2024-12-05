package repositories

import (
	"database/sql"
	"time"

	"github.com/filipe1309/agl-go-driver/internal/users"
)

type UserReadRepository interface {
	Login(string, string) *sql.Row
	ReadDB(int64) *sql.Row
	ReadAllDB() *sql.Rows
}

type UserWriteRepository interface {
	InsertDB(*users.User) (int64, error)
	UpdateDB(int64, *users.User) (int64, error)
	SoftDeleteDB(int64) error
}

type UserReadWriteRepository interface {
	UserReadRepository
	UserWriteRepository
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) InsertDB(user *users.User) (id int64, err error) {
	stmt := `INSERT INTO users (name, login, password, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`

	err = ur.db.QueryRow(stmt, user.Name, user.Login, user.Password, user.UpdatedAt).Scan(&id)
	if err != nil {
		return -1, err
	}

	return
}

func (ur *UserRepository) UpdateDB(id int64, user *users.User) (int64, error) {
	user.UpdatedAt = time.Now()
	stmt := `UPDATE users SET name = $1, updated_at = $2, last_login = $3 WHERE id = $4`

	result, err := ur.db.Exec(stmt, user.Name, user.UpdatedAt, user.LastLogin, id)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (ur *UserRepository) SoftDeleteDB(id int64) error {
	stmt := `UPDATE users SET updated_at = $1, deleted = TRUE WHERE id = $2`

	_, err := ur.db.Exec(stmt, time.Now(), id)

	return err
}

func (ur *UserRepository) ReadDB(id int64) *sql.Row {
	stmt := `SELECT * FROM users WHERE id = $1`
	row := ur.db.QueryRow(stmt, id)
	return row
}

func (ur *UserRepository) ReadAllDB() *sql.Rows {
	stmt := `SELECT * FROM users WHERE deleted = FALSE`

	rows, err := ur.db.Query(stmt)
	if err != nil {
		return nil
	}
	defer rows.Close()

	return rows
}

func (ur *UserRepository) Login(login, password string) *sql.Row {
	stmt := `SELECT * FROM users WHERE login = $1 AND password = $2`
	return ur.db.QueryRow(stmt, login, password)
}
