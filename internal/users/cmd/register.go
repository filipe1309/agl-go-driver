package cmd

import "github.com/spf13/cobra"

func Register(c *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage users",
	}

	cmd.AddCommand(create())
	cmd.AddCommand(update())
	cmd.AddCommand(delete())
	cmd.AddCommand(list())

	c.AddCommand(cmd)
}
