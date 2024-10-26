//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	// Используем unshare для создания нового PID namespace
	cmd := exec.Command("unshare", "--fork", "--pid", "--mount", "--user", "--map-root-user", "--mount-proc", "bash")

	// Привязываем стандартный ввод/вывод к текущему процессу
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Запуск команды
	if err := cmd.Start(); err != nil {
		fmt.Printf("Ошибка при запуске bash в новом namespace: %v\n", err)
		os.Exit(1)
	}

	// Ждем завершение команды
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Printf("bash завершился с ошибкой: %v\n", err)
		}
	}()

	// Обработка сигналов
	exitSignals := make(chan os.Signal, 1)
	signal.Notify(exitSignals, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнал завершения
	<-exitSignals
	fmt.Println("\nПолучен сигнал завершения, завершаем bash.")
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("Ошибка при завершении bash: %v\n", err)
	}
}
