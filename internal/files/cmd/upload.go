package cmd

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/filipe1309/agl-go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func upload() *cobra.Command {
	var (
		filename string
		folderID int32
	)

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a file",
		Run: func(cmd *cobra.Command, args []string) {
			if filename == "" {
				log.Fatal("Please provide a name")
			}

			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			var body bytes.Buffer
			mw := multipart.NewWriter(&body)
			fw, err := mw.CreateFormFile("file", filepath.Base(file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(fw, file)
			if err != nil {
				log.Fatal(err)
			}

			if folderID > 0 {
				w, err := mw.CreateFormField("folder_id")
				if err != nil {
					log.Fatal(err)
				}
				_, err = w.Write([]byte(strconv.Itoa(int(folderID))))
				if err != nil {
					log.Fatal(err)
				}
			}

			err = mw.Close()
			if err != nil {
				log.Fatal(err)
			}

			headers := map[string]string{
				"Content-Type": mw.FormDataContentType(),
			}

			_, err = requests.Post("/files", &body, headers, true)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("File %s uploaded", filename)
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "", "File name")
	cmd.Flags().Int32VarP(&folderID, "folder", "d", 0, "Folder ID")

	return cmd
}
