package cmd

import (
	"log"

	"github.com/filipe1309/agl-go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func authenticate() *cobra.Command {
	var (
		username string
		password string
	)

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate user in the API",
		Run: func(cmd *cobra.Command, args []string) {
			if username == "" || password == "" {
				log.Fatal("Username and password are required")
			}

			err := requests.Auth("/auth", username, password)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&username, "user", "u", "", "Username")
	cmd.Flags().StringVarP(&password, "pass", "p", "", "Password")

	return cmd
}
