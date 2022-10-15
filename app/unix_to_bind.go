package app

import (
	"github.com/spf13/cobra"
	"log"
	"net"
)

func (a *App) UnixSocket2Bind() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unix2bind",
		Short:   "unix socket to bind",
		Aliases: []string{"u2b"},
		Args:    cobra.ExactArgs(2),
		Example: "./red u2b /var/run/docker.sock 0.0.0.0:2345",
	}

	cmd.RunE = a.UnixSocket2BindHandler
	return cmd
}

func (a *App) UnixSocket2BindHandler(_ *cobra.Command, args []string) error {
	leftAddr := args[0]
	rightAddr := args[1]
	log.Println("UnixSocket to Connect:", leftAddr, "<----->", rightAddr)
	if right, err := net.Listen("tcp", rightAddr); err == nil {
		for {
			if rightConn, err := right.Accept(); err == nil {
				if leftConn, err := net.Dial("unix", leftAddr); err == nil {
					go func(left, right net.Conn) {
						if err := a.Forward(left, right); err != nil {
							log.Println(err)
						}
					}(rightConn, leftConn)
				} else {
					log.Printf("Error: connect %s failed.(%s)\n", leftAddr, err)
				}
			} else {
				log.Println(err)
			}
		}
	} else {
		log.Printf("Error: bind %s failed.(%s)\n", rightAddr, err)
	}
	return nil
}
