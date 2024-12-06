package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/filipe1309/agl-go-driver/internal/users"
	"github.com/filipe1309/agl-go-driver/pkg/requests"
	userspb "github.com/filipe1309/agl-go-driver/proto/v1/users"
	"github.com/spf13/cobra"
)

func create() *cobra.Command {
	var (
		name     string
		login    string
		password string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" || login == "" || password == "" {
				log.Fatal("Please provide a name, login and password")
			}

			mode := cmd.Flag("mode").Value.String()
			switch mode {
			case "http":
				createWithHTTP(name, login, password)
			case "grpc":
				createWithGRPC(name, login, password)
			default:
				log.Fatalf("Mode %s not supported", mode)
			}

			log.Printf("User created")
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "User name")
	cmd.Flags().StringVarP(&login, "login", "l", "", "User login")
	cmd.Flags().StringVarP(&password, "pass", "p", "", "User password")

	return cmd
}

func createWithHTTP(name, login, password string) {
	user := users.User{
		Name:     name,
		Login:    login,
		Password: password,
	}
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(user)
	if err != nil {
		log.Fatal(err)
	}

	_, err = requests.Post("/users", &body, nil, false)
	if err != nil {
		log.Fatal(err)
	}
}

func createWithGRPC(name, login, password string) {
	user := &userspb.UserRequest{
		Name:     name,
		Login:    login,
		Password: password,
	}

	conn := requests.GetGRPCConn()
	defer conn.Close()

	client := userspb.NewUserServiceClient(conn)

	resp, err := client.Create(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("User created with id: %v", resp.User.Id)
}
