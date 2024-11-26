package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/filipe1309/agl-go-driver/internal/folders"
	"github.com/filipe1309/agl-go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func list() *cobra.Command {
	var id int32

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all folders",
		Run: func(cmd *cobra.Command, args []string) {
			path := "/folders"

			if id > 0 {
				path = fmt.Sprintf("%s/%d", path, id)
			}

			data, err := requests.AuthenticatedGet(path)
			if err != nil {
				log.Fatal(err)
			}

			var fc folders.FolderContent
			err = json.Unmarshal(data, &fc)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(fc.Folder.Name)
			log.Println("====================")
			for _, fr := range fc.Content {
				log.Println(fr.ID, " - ", fr.Type, " - ", fr.Name)
			}
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "Folder ID")

	return cmd
}
