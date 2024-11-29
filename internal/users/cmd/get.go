package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/filipe1309/agl-go-driver/internal/users"
	"github.com/filipe1309/agl-go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func get() *cobra.Command {
	var id int32

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get user by ID",
		Run: func(cmd *cobra.Command, args []string) {
			if id <= 0 {
				log.Fatal("ID is required")
			}
			path := "/users"
			path = fmt.Sprintf("%s/%d", path, id)

			data, err := requests.AuthenticatedGet(path)
			if err != nil {
				log.Fatal(err)
			}

			var u users.User
			err = json.Unmarshal(data, &u)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(u.Name)
			log.Println(u.Login)
			log.Println(u.LastLogin)
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "Folder ID")

	return cmd
}
