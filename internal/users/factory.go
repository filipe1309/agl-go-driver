package users

import (
	"database/sql"
	"log"
)

type scan interface {
	Scan(dest ...any) error
}

func restore(row scan) (*User, error) {
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.Deleted)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func RestoreOne(row *sql.Row) (*User, error) {
	return restore(row)
}

func RestoreAll(rows *sql.Rows) ([]User, error) {
	users := make([]User, 0)
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
