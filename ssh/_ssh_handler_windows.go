package ssh

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unicode/utf8"
)

func sshHandler(sess ssh.Session) {
	_, winCh, isPty := sess.Pty()
	if isPty {
		term := terminal.NewTerminal(sess, "")

		go func() {
			for win := range winCh {
				if err := term.SetSize(win.Width, win.Height); err != nil {
					log.Println(err)
				}
			}
		}()

		for {
			curDir, err := os.Getwd()
			if err != nil {
				curDir = err.Error()
			}
			term.SetPrompt(string(term.Escape.Cyan) + fmt.Sprintf("[%s]> ", curDir) + string(term.Escape.Reset))
			line, err := term.ReadLine()
			if err != nil {
				log.Println(err)
				break
			}
			if line == "" {
				continue
			}
			if strings.HasPrefix(strings.ToLower(line), "cd ") {
				if err := os.Chdir(strings.TrimSpace(line[2:])); err != nil {
					_, _ = term.Write([]byte(fmt.Sprintf("%s\n", err.Error())))
				}
				continue
			}

			if line == "exit" {
				break
			}
			if _, err := term.Write(runCommand(line)); err != nil {
				log.Println(err)
				break
			}
		}
	} else {
		remoteCommand := strings.Join(sess.Command(), " ")
		result := runCommand(remoteCommand)
		_, _ = sess.Write(result)
	}
	_ = sess.Exit(0)
}

func convert(bs []byte) []byte {
	if !utf8.Valid(bs) {
		for _, charset := range simplifiedchinese.All {
			if result, err := charset.NewDecoder().Bytes(bs); err == nil {
				return result
			}
		}
	}
	return bs
}

func runCommand(line string) []byte {
	cmd := exec.Command("cmd.exe", "/C", line)
	if pwd, err := os.Getwd(); err == nil {
		cmd.Dir = pwd
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if result, err := cmd.CombinedOutput(); err == nil {
		return convert(result)
	} else {
		return convert(result)
	}
}
