package app

import (
	"github.com/spf13/cobra"
	"log"
	"net"
)

func (a *App) Bind2Bind() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bind2bind",
		Short:   "bind to bind",
		Aliases: []string{"b2b"},
		Args:    cobra.ExactArgs(2),
		Example: "./red b2b 0.0.0.0:1234 0.0.0.0:5678",
	}

	cmd.RunE = a.Bind2BindHandler
	return cmd
}

func (a *App) Bind2BindHandler(_ *cobra.Command, args []string) error {
	leftAddr := args[0]
	rightAddr := args[1]
	log.Println("Bind to Bind:", leftAddr, "<----->", rightAddr)
	if left, err := net.Listen("tcp", leftAddr); err == nil {
		if right, err := net.Listen("tcp", rightAddr); err == nil {
			for {
				if leftConn, err := left.Accept(); err == nil {
					log.Println("Accept:", leftAddr)
					if rightConn, err := right.Accept(); err == nil {
						log.Println("Accept:", rightAddr)
						go func(left, right net.Conn) {
							if err := a.Forward(left, right); err != nil {
								log.Println(err)
							}
						}(leftConn, rightConn)
					} else {
						log.Println(err)
					}
				} else {
					log.Println(err)
				}
			}
		} else {
			log.Printf("Error: bind %s failed.(%s)\n", rightAddr, err)
		}
	} else {
		log.Printf("Error: bind %s failed.(%s)\n", leftAddr, err)
	}
	return nil
}
