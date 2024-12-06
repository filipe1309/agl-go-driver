package main

import (
	"log"

	"github.com/filipe1309/agl-go-driver/cmd/api/server"
	"github.com/spf13/cobra"
)

func main() {
	var mode string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the API with the specified mode (http, grpc)",
		Run: func(cmd *cobra.Command, args []string) {
			if mode == "" {
				log.Fatal("Mode is required")
			}

			switch mode {
			case "http":
				server.RunHTTPServer()
			case "grpc":
				server.RunGRPCServer()
			default:
				log.Fatalf("Mode %s not supported", mode)
			}
		},
	}

	cmd.Flags().StringVarP(&mode, "mode", "m", "http", "Mode of operation (http, grpc)")

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
