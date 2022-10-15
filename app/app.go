package app

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

type App struct {
	rootCmd *cobra.Command

	ctx    context.Context
	cancel context.CancelFunc

	background bool
}

func (a *App) Init() error {
	// 后台运行模式
	a.rootCmd.PersistentFlags().BoolVarP(&a.background, "background", "b", false, "background mode")

	a.rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		// 判断是否需要进入后台模式
		if a.background && os.Getppid() != 1 {
			// Windows 不支持后台模式
			if runtime.GOOS == "windows" {
				return fmt.Errorf("windows is not support background mode")
			}
			_ = exec.Command(os.Args[0], os.Args[1:]...).Start()
			os.Exit(0)
		}
		return nil
	}
	return nil
}

func (a *App) Run() error {
	defer a.cancel()

	// 初始化
	if err := a.Init(); err != nil {
		return err
	}

	// 注册子命令
	a.registerCommands()
	a.rootCmd.CompletionOptions.HiddenDefaultCmd = true
	return a.rootCmd.Execute()
}

func (a *App) registerCommands() {
	// I/O(网络/Unix套接字)流量流转相关
	a.rootCmd.AddCommand(a.Bind2Bind())          // b2b
	a.rootCmd.AddCommand(a.Bind2Connect())       // b2c
	a.rootCmd.AddCommand(a.Connect2Connect())    // c2c
	a.rootCmd.AddCommand(a.UnixSocket2Bind())    // u2b
	a.rootCmd.AddCommand(a.UnixSocket2Connect()) // u2c

	// 代理相关
	//a.rootCmd.AddCommand(a.Socks2Connect())
	a.rootCmd.AddCommand(a.Socks5()) // s5

	// SSH隧道相关
	a.rootCmd.AddCommand(a.BindSSH())    // bs
	a.rootCmd.AddCommand(a.ReverseSSH()) // revssh / rev / rs

	// 端口扫描
	a.rootCmd.AddCommand(a.PortScan()) // ps / scan

	// HTTP静态文件服务
	a.rootCmd.AddCommand(a.HTTPStaticServer())
	// 其他
}

// Pipe 流量流转实现
func Pipe(rwc1, rwc2 io.ReadWriteCloser) error {
	errCh := make(chan error, 2)
	pipe := func(dst io.WriteCloser, src io.ReadCloser, wg *sync.WaitGroup) error {
		defer func() {
			_ = dst.Close()
			wg.Done()
		}()
		_, err := io.Copy(dst, src)
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		errCh <- pipe(rwc1, rwc2, wg)
	}()
	go func() {
		errCh <- pipe(rwc2, rwc1, wg)
	}()

	wg.Wait()

	if err := <-errCh; err == nil || err == io.EOF {
		return nil
	} else {
		return err
	}
}

func (a *App) Forward(conn1, conn2 net.Conn) error {
	log.Printf("[(%s)==(%s)]<===>[(%s)==(%s)]\n", conn2.LocalAddr(), conn2.RemoteAddr(), conn1.LocalAddr(), conn1.RemoteAddr())
	defer func() {
		log.Printf("[(%s)==(%s)]<=X=>[(%s)==(%s)]\n", conn2.LocalAddr(), conn2.RemoteAddr(), conn1.LocalAddr(), conn1.RemoteAddr())
	}()
	return Pipe(conn1, conn2)
}

func New() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		rootCmd: &cobra.Command{
			Use: "redirector",
		},
		ctx:    ctx,
		cancel: cancel,
	}
}
