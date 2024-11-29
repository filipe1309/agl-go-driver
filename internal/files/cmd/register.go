package cmd

import "github.com/spf13/cobra"

func Register(c *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Manage files",
	}

	cmd.AddCommand(upload())
	cmd.AddCommand(update())
	cmd.AddCommand(delete())

	c.AddCommand(cmd)
}
