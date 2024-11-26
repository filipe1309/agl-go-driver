package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/filipe1309/agl-go-driver/internal/folders"
	"github.com/filipe1309/agl-go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func update() *cobra.Command {
	var id int32
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a folder name",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" || id == 0 {
				log.Fatal("Please provide a name and an ID")
			}

			folder := folders.Folder{Name: name}
			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(folder)
			if err != nil {
				log.Fatal(err)
			}

			path := fmt.Sprintf("/folders/%d", id)
			_, err = requests.AuthenticatedPut(path, &body)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Folder %d updated with name %s", id, name)
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "Folder ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Folder name")

	return cmd
}
