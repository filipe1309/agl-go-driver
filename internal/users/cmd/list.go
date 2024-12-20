package cmd

import (
	"encoding/json"
	"log"

	"github.com/filipe1309/agl-go-driver/internal/users"
	"github.com/filipe1309/agl-go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func list() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Run: func(cmd *cobra.Command, args []string) {
			path := "/users"
			data, err := requests.AuthenticatedGet(path)
			if err != nil {
				log.Fatal(err)
			}

			var us []users.User
			err = json.Unmarshal(data, &us)
			if err != nil {
				log.Fatal(err)
			}

			for _, u := range us {
				log.Println(u.Name, u.Login, u.LastLogin)
			}
		},
	}

	return cmd
}
