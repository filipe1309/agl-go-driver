package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/filipe1309/agl-go-driver/internal/files"
	"github.com/filipe1309/agl-go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func update() *cobra.Command {
	var id int32
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update file name",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" || id <= 0 {
				log.Fatal("Please provide a name and an ID")
			}

			file := files.File{Name: name}
			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(file)
			if err != nil {
				log.Fatal(err)
			}

			path := fmt.Sprintf("/files/%d", id)
			_, err = requests.AuthenticatedPut(path, &body)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("File %d updated with name %s", id, name)
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "Folder ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Folder name")

	return cmd
}
