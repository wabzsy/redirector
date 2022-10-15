package ssh

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/pkg/sftp"
	"io"
	"log"
	"redirector/utils"
	"runtime"
)

func NewSSHServer(password string) (*ssh.Server, error) {
	return newSSHServer(sshHandler,
		SetHostKey(),
		SetLoginOption(password),
		SetServerVersion(),
		SetPortForwardingHandler(),
		SetSftpHandler(),
	)
}

func newSSHServer(handler ssh.Handler, options ...ssh.Option) (*ssh.Server, error) {
	srv := &ssh.Server{Handler: handler}

	// ensureHandler
	if srv.RequestHandlers == nil {
		srv.RequestHandlers = map[string]ssh.RequestHandler{}
		for k, v := range ssh.DefaultRequestHandlers {
			srv.RequestHandlers[k] = v
		}
	}
	if srv.ChannelHandlers == nil {
		srv.ChannelHandlers = map[string]ssh.ChannelHandler{}
		for k, v := range ssh.DefaultChannelHandlers {
			srv.ChannelHandlers[k] = v
		}
	}
	if srv.SubsystemHandlers == nil {
		srv.SubsystemHandlers = map[string]ssh.SubsystemHandler{}
		for k, v := range ssh.DefaultSubsystemHandlers {
			srv.SubsystemHandlers[k] = v
		}
	}

	for _, option := range options {
		if err := srv.SetOption(option); err != nil {
			return nil, err
		}
	}

	// ensureSigner
	if len(srv.HostSigners) == 0 {
		signer, err := utils.GenerateSigner()
		if err != nil {
			return nil, err
		}
		srv.HostSigners = append(srv.HostSigners, signer)
	}
	return srv, nil
}

func SetHostKey() ssh.Option {
	key := `-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBALkdaxicFBMZYxzt
vKom45LvVSne5Ag6M+kXQk9ksN1k+EqZj7h/3ctxDEjOVdqv9CgwPLgM1iuEmqYI
9lJNlhjEj8ZLbg8fSymjbFoAHdQuErMEcAAGYXAtv+7zj75AZFvlC4pFBhfrUC1o
hcdbtmGtepP7nk3G4vpiCVdid8CjAgMBAAECgYAtIws2GPicH5iXOTDDnG/pKApw
BzU6/FYkA9PbYAXwNeqE5iSxLBx8urfwGL++byDtm7Vye07NlavPyGencdujDB4b
ldawX8FKTh1CFFrvRBaN2lto0r0ejllNAj4MvBslGeXwtqvB3NXV0gul55tbVgLO
s+nsxSVW8ALmgt0c4QJBAOUPDVmy4eXpUmLV/sIsOCuiXZNjKJcMhXT39Tscrshd
R6jo0UAJ/quwUygxqM4kevt69dQ/6hlxWCfTo3M4tg0CQQDO4zfjRtX4062cGcgW
lRV0/CAcw71Be7qKxiCf25dpjCdxNZWjMORfeiGCoMKLwWbE/vcueLcf69VHD3iB
CNVvAkBeH4tK2pi80t2Jw4mF3InQVE3DbLGXMAv+/o0El0qzBrGVlOW3POQrRK9H
CvDklFT81ZACgaj+f3bMFslJZXpZAkBux1PhqshgGFhZwaRWEzYOEgLP5C+upKXa
MQS/FEIbDiUAhYS+gSuHxEm1PIdvdfuleDC6/YBw40KsbihET4qZAkAC0nu/Gkly
GMbBKpfRyFxg31hgY/yQIMYe7XJ3lCmqv14J8o9Gyf++o5FtP/L/Smjr0V4E8lLP
BZGEhvLIryFk
-----END PRIVATE KEY-----
`
	return ssh.HostKeyPEM([]byte(key))
}

func SetServerVersion() ssh.Option {
	return func(srv *ssh.Server) error {
		srv.Version = fmt.Sprintf("OpenSSL_7.6-%s-%s", runtime.GOOS, runtime.GOARCH)
		return nil
	}
}

func SetSftpHandler() ssh.Option {
	return func(srv *ssh.Server) error {
		srv.SubsystemHandlers["sftp"] = SftpHandler
		return nil
	}
}

func SftpHandler(sess ssh.Session) {
	server, err := sftp.NewServer(sess)
	if err != nil {
		log.Printf("sftp server init error: %s\n", err)
		return
	}
	if err := server.Serve(); err == io.EOF {
		_ = server.Close()
		log.Println("sftp client exited session.")
	} else if err != nil {
		log.Println("sftp server completed with error:", err)
	}
}

func SetLoginOption(password string) ssh.Option {
	return ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool {
		return pass == password
	})
}

func SetPortForwardingHandler() ssh.Option {
	return func(srv *ssh.Server) error {
		forwardHandler := &ssh.ForwardedTCPHandler{}
		srv.RequestHandlers["tcpip-forward"] = forwardHandler.HandleSSHRequest
		srv.RequestHandlers["cancel-tcpip-forward"] = forwardHandler.HandleSSHRequest
		srv.ReversePortForwardingCallback = func(ctx ssh.Context, host string, port uint32) bool {
			// -R
			//log.Println("attempt to bind", host, port, "granted")
			return true
		}
		srv.ChannelHandlers["direct-tcpip"] = ssh.DirectTCPIPHandler
		srv.LocalPortForwardingCallback = func(ctx ssh.Context, dhost string, dport uint32) bool {
			// -L
			//log.Println("Accepted forward", dhost, dport)
			return true
		}
		return nil
	}
}
