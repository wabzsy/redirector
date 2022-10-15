package app

import (
	"github.com/spf13/cobra"
	"log"
	"redirector/ssh"
)

func (a *App) BindSSH() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bind-ssh",
		Short:   "run ssh server",
		Aliases: []string{"bs"},
		Args:    cobra.ExactArgs(1),
		Example: "./red bs 0.0.0.0:2345",
	}

	cmd.Flags().StringP("password", "p", "123456", "ssh service password")

	cmd.RunE = a.BindSSHHandler
	return cmd
}

func (a *App) BindSSHHandler(cmd *cobra.Command, args []string) error {
	addr := args[0]

	password, err := cmd.Flags().GetString("password")

	if err != nil {
		return err
	}

	log.Println("Serving SSH service on:", addr, "password:", password)

	server, err := ssh.NewSSHServer(password)
	if err != nil {
		return err
	}
	server.Addr = addr
	return server.ListenAndServe()
}
