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

			mode := cmd.Flag("mode").Value.String()
			switch mode {
			case "http":
				authWithHTTP(username, password)
			case "grpc":
				authWithGRPC(username, password)
			default:
				log.Fatalf("Mode %s not supported", mode)
			}
		},
	}

	cmd.Flags().StringVarP(&username, "user", "u", "", "Username")
	cmd.Flags().StringVarP(&password, "pass", "p", "", "Password")

	return cmd
}

func authWithHTTP(username, password string) {
	err := requests.HTTPAuth("/auth", username, password)
	if err != nil {
		log.Fatal(err)
	}
}

func authWithGRPC(username, password string) {
	err := requests.GRPCAuth(username, password)
	if err != nil {
		log.Fatal(err)
	}
}
