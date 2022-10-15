package app

import (
	"github.com/spf13/cobra"
	"log"
	"net"
	"time"
)

func (a *App) UnixSocket2Connect() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unix2connect",
		Short:   "unix socket to connect",
		Aliases: []string{"u2c"},
		Args:    cobra.ExactArgs(2),
		Example: "./red u2c /var/run/docker.sock 1.2.3.4:2345",
	}

	cmd.RunE = a.UnixSocket2ConnectHandler
	return cmd
}

func (a *App) UnixSocket2ConnectHandler(_ *cobra.Command, args []string) error {
	leftAddr := args[0]
	rightAddr := args[1]
	log.Println("UnixSocket to Connect:", leftAddr, "<----->", rightAddr)

	// 限制最大20个连接
	maxConn := make(chan struct{}, 20)

	for {
		if leftConn, err := net.Dial("unix", leftAddr); err == nil {
			if rightConn, err := net.Dial("tcp", rightAddr); err == nil {
				maxConn <- struct{}{}
				go func(left, right net.Conn) {
					defer func() {
						<-maxConn
					}()
					if err := a.Forward(left, right); err != nil {
						log.Println(err)
					}
				}(rightConn, leftConn)
			} else {
				log.Printf("Error: connect %s failed.(%s)\n", rightAddr, err)
				_ = leftConn.Close()
				time.Sleep(time.Second * 3)
			}
		} else {
			log.Printf("Error: connect %s failed.(%s)\n", leftAddr, err)
			time.Sleep(time.Second * 3)
			continue
		}
	}
}
