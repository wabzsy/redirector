package ssh

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os/exec"
	"redirector/utils"
	"syscall"
)

func DefaultCommand() string {
	cmd := "C:\\Windows\\System32\\cmd.exe"
	ps := "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
	if utils.FileExists(cmd) {
		return cmd
	} else {
		return ps
	}
}

func sshHandler(sess ssh.Session) {
	cmd := GetCommand(sess)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, winCh, isPty := sess.Pty()
	if isPty {
		term := terminal.NewTerminal(sess, "")
		_, _ = term.Write([]byte("Tips: chcp 65001\n"))
		go func() {
			for win := range winCh {
				if err := term.SetSize(win.Width, win.Height); err != nil {
					log.Println(err)
				}
			}
		}()

		cmd.Stderr = sess
		cmd.Stdout = sess
		stdin, err := cmd.StdinPipe()
		if err != nil {
			writeError(sess, err)
			return
		}

		if err := cmd.Start(); err != nil {
			writeError(sess, err)
			return
		}

		go func() {
			for {
				line, err := term.ReadLine()
				if err != nil {
					if err == io.EOF {
						if cmd.Process != nil {
							kill := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(cmd.Process.Pid))
							kill.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
							_ = kill.Run()
						}
					} else {
						writeError(sess, err)
					}
					break
				}
				if _, err = stdin.Write([]byte(line + "\n")); err != nil {
					writeError(sess, err)
					break
				}
			}
		}()
		if err := cmd.Wait(); err != nil {
			writeError(sess, err)
			return
		}
	} else {
		if result, err := cmd.CombinedOutput(); err == nil {
			_, _ = sess.Write(result)
		} else {
			_, _ = sess.Write([]byte(err.Error()))
			_ = sess.Exit(1)
		}
	}

	_ = sess.Exit(0)
}
