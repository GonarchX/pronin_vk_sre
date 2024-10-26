package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	var overall [][]byte
	var m runtime.MemStats

	go func() {
		for {
			a := make([]byte, 0, 1*1024*1024) // 1M
			overall = append(overall, a)

			runtime.ReadMemStats(&m)
			fmt.Printf("Всего выделенной памяти: %v МБ\n", m.Alloc/(1024*1024))
			time.Sleep(1 * time.Second)
		}
	}()

	<-exit
}
