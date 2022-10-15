package main

import (
	"os/signal"
	"syscall"
)

func init() {
	signal.Ignore(syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
}
