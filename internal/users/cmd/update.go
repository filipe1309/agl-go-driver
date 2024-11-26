package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/filipe1309/agl-go-driver/internal/users"
	"github.com/filipe1309/agl-go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func update() *cobra.Command {
	var id int32
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a user name",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" || id == 0 {
				log.Fatal("Please provide a name and an ID")
			}

			user := users.User{Name: name}
			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(user)
			if err != nil {
				log.Fatal(err)
			}

			path := fmt.Sprintf("/users/%d", id)
			_, err = requests.AuthenticatedPut(path, &body)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("User %d updated with name %s", id, name)
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "Folder ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Folder name")

	return cmd
}
