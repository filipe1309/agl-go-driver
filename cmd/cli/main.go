package main

import (
	"log"

	authCmd "github.com/filipe1309/agl-go-driver/internal/auth/cmd"
	filesCmd "github.com/filipe1309/agl-go-driver/internal/files/cmd"
	foldersCmd "github.com/filipe1309/agl-go-driver/internal/folders/cmd"
	usersCmd "github.com/filipe1309/agl-go-driver/internal/users/cmd"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{}

func main() {
	var mode string

	RootCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "http", "Mode of operation (http, grpc)")

	authCmd.Register(RootCmd)
	filesCmd.Register(RootCmd)
	foldersCmd.Register(RootCmd)
	usersCmd.Register(RootCmd)

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
