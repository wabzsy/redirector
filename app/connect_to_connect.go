package app

import (
	"github.com/spf13/cobra"
	"log"
	"net"
	"time"
)

func (a *App) Connect2Connect() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connect2connect",
		Short:   "connect to connect",
		Aliases: []string{"c2c"},
		Args:    cobra.ExactArgs(2),
		Example: "./red c2c 1.2.3.4:8080 5.6.7.8:8080",
	}

	cmd.RunE = a.Connect2ConnectHandler
	return cmd
}

func (a *App) Connect2ConnectHandler(_ *cobra.Command, args []string) error {
	leftAddr := args[0]
	rightAddr := args[1]
	log.Println("Connect to Connect:", leftAddr, "<----->", rightAddr)
	for {
		if leftConn, err := net.Dial("tcp", leftAddr); err == nil {
			if rightConn, err := net.Dial("tcp", rightAddr); err == nil {
				go func(left, right net.Conn) {
					if err := a.Forward(left, right); err != nil {
						log.Println(err)
					}
				}(leftConn, rightConn)
			} else {
				log.Printf("Error: connect %s failed.(%s)\n", rightAddr, err)
				_ = leftConn.Close()
			}
		} else {
			log.Printf("Error: connect %s failed.(%s)\n", leftAddr, err)
		}
		time.Sleep(time.Second)
	}
}
