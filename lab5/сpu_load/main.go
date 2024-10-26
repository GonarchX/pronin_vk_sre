package main

import (
	"os"
	"os/signal"
	"syscall"
	//"time"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-exit:
			return
		default:
		}
	}
}
