package app

import (
	"github.com/armon/go-socks5"
	"github.com/spf13/cobra"
	"log"
)

func (a *App) Socks5() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "socks5",
		Short:   "socks5 server",
		Aliases: []string{"s5"},
		Args:    cobra.ExactArgs(1),
		Example: "./red socks5 0.0.0.0:1080",
	}

	cmd.RunE = a.Socks5Handler
	return cmd
}

func (a *App) Socks5Handler(_ *cobra.Command, args []string) error {
	addr := args[0]
	log.Println("Serving Socks5 on:", addr)
	if server, err := socks5.New(&socks5.Config{}); err == nil {
		if err = server.ListenAndServe("tcp", addr); err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
	return nil
}
