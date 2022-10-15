//go:build !windows
// +build !windows

package ssh

import (
	"fmt"
	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"io"
	"os"
	"redirector/errors"
	"redirector/utils"
	"sync"
	"syscall"
	"unsafe"
)

func DefaultCommand() string {
	zsh := "/bin/zsh"
	bash := "/bin/bash"
	sh := "/bin/sh"
	if utils.FileExists(zsh) {
		return zsh
	} else if utils.FileExists(bash) {
		return bash
	} else {
		return sh
	}
}

func sshHandler(sess ssh.Session) {
	cmd := GetCommand(sess)

	cmd.Env = []string{
		"HISTFILE=/dev/null",
		"LC_ALL=en_US.UTF-8",
		"LANG=en_US.UTF-8",
		//"HISTSIZE=0",
		//"HISTFILESIZE=0",
		"PATH=/usr/local/bin:/usr/local/sbin:/usr/bin:/bin:/usr/sbin:/sbin",
	}

	ptyReq, winCh, isPty := sess.Pty()

	if isPty {
		cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
		f, err := pty.Start(cmd)
		if err != nil {
			writeError(sess, errors.New("PTY start failed.\n"))
			return
		}

		go func() {
			for win := range winCh {
				setWinSize(f, win.Width, win.Height)
			}
		}()

		doneCh := make(chan bool, 1)
		var once sync.Once

		done := func() {
			_ = cmd.Wait()
			_ = f.Close()
			doneCh <- true
		}

		go func() {
			_, _ = io.Copy(f, sess) // stdin
			once.Do(done)
		}()
		go func() {
			_, _ = io.Copy(sess, f) // stdout
			once.Do(done)
		}()

		<-doneCh
	} else {
		if result, err := cmd.CombinedOutput(); err == nil {
			_, _ = sess.Write(result)
		} else {
			writeError(sess, err)
			return
		}
	}
	_ = sess.Exit(0)
}

func setWinSize(f *os.File, w, h int) {
	_, _, _ = syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(syscall.TIOCSWINSZ),
		uintptr(
			unsafe.Pointer(
				&struct {
					h, w, x, y uint16
				}{
					uint16(h),
					uint16(w),
					uint16(0),
					uint16(0),
				},
			),
		),
	)
}
