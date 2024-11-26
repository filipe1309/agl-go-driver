package cmd

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/filipe1309/agl-go-driver/internal/folders"
	"github.com/filipe1309/agl-go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func create() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new folder",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" {
				log.Fatal("Please provide a name")
			}

			folder := folders.Folder{Name: name}
			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(folder)
			if err != nil {
				log.Fatal(err)
			}

			_, err = requests.Post("/folders", &body, true)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Folder %s created", name)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Folder name")

	return cmd
}
