package cmd

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/filipe1309/agl-go-driver/internal/users"
	"github.com/filipe1309/agl-go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func create() *cobra.Command {
	var (
		name string
		login string
		password string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" || login == "" || password == "" {
				log.Fatal("Please provide a name, login and password")
			}

			user := users.User{
				Name: name,
				Login: login,
				Password: password,
			}
			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(user)
			if err != nil {
				log.Fatal(err)
			}

			_, err = requests.Post("/users", &body, false)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("User %s created", name)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "User name")
	cmd.Flags().StringVarP(&login, "login", "l", "", "User login")
	cmd.Flags().StringVarP(&password, "pass", "p", "", "User password")

	return cmd
}
