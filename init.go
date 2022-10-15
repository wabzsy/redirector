//go:build !windows

package main

import (
	"os/signal"
	"syscall"
)

func init() {
	signal.Ignore(syscall.SIGHUP, syscall.SIGTSTP, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGQUIT)
}
