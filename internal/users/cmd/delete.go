package cmd

import (
	"fmt"
	"log"

	"github.com/filipe1309/agl-go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func delete() *cobra.Command {
	var id int32

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a user",
		Run: func(cmd *cobra.Command, args []string) {
			if id <= 0 {
				log.Fatal("Please provide an ID")
			}
			path := "/users"
			path = fmt.Sprintf("%s/%d", path, id)

			err := requests.AuthenticatedDelete(path)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("User %d deleted", id)
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "User ID")

	return cmd
}
