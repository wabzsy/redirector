package app

import (
	"github.com/spf13/cobra"
	"log"
	"net"
	"time"
)

func (a *App) Bind2Connect() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bind2connect",
		Short:   "bind to connect",
		Aliases: []string{"b2c"},
		Args:    cobra.ExactArgs(2),
		Example: "./red b2c 0.0.0.0:1234 1.2.3.4:1234",
	}

	cmd.RunE = a.Bind2ConnectHandler
	return cmd
}

func (a *App) Bind2ConnectHandler(_ *cobra.Command, args []string) error {
	leftAddr := args[0]
	rightAddr := args[1]
	log.Println("Bind to Connect:", leftAddr, "<----->", rightAddr)
	if left, err := net.Listen("tcp", leftAddr); err == nil {
		for {
			if leftConn, err := left.Accept(); err == nil {
				if rightConn, err := net.Dial("tcp", rightAddr); err == nil {
					go func(left, right net.Conn) {
						if err := a.Forward(left, right); err != nil {
							log.Println(err)
						}
					}(leftConn, rightConn)
				} else {
					log.Printf("Error: connect %s failed.(%s)\n", rightAddr, err)
					_ = leftConn.Close()
					time.Sleep(time.Second * 3)
				}
			} else {
				log.Println(err)
			}
		}
	} else {
		log.Printf("Error: bind %s failed.(%s)\n", leftAddr, err)
	}
	return nil
}
