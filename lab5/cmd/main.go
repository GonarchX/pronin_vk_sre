package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	memoryLimitFile = "memory.max"
	cpuLimitFile    = "cpu.max"
	procsFile       = "cgroup.procs"
)

func createCgroup(cgroupDirectory string) error {
	if err := os.MkdirAll(cgroupDirectory, 0755); err != nil {
		return fmt.Errorf("ошибка при создании директории %s: %w", cgroupDirectory, err)
	}

	return nil
}

func setLimits(cgroupDirectory, memoryLimit, cpuLimit string) error {
	fileName := filepath.Join(cgroupDirectory, memoryLimitFile)

	if err := os.WriteFile(fileName, []byte(memoryLimit), 0644); err != nil {
		return fmt.Errorf("ошибка при задании лимитов по памяти: %w", err)
	}

	fileName = filepath.Join(cgroupDirectory, cpuLimitFile)

	if err := os.WriteFile(fileName, []byte(cpuLimit), 0644); err != nil {
		return fmt.Errorf("ошибка при задании лимитов по процессорному времени: %w", err)
	}

	return nil
}

func bindCurrentProcessToCgroup(cgroupDirectory string) error {
	fileName := filepath.Join(cgroupDirectory, procsFile)
	curPid := []byte(strconv.Itoa(os.Getpid()))

	if err := os.WriteFile(fileName, curPid, 0644); err != nil {
		return fmt.Errorf("ошибка при привязке процесса к cgroup: %w", err)
	}

	return nil
}

func runCommandInNamespace(command string, args []string) error {
	cmd := exec.Command(command, args...)

	// Привязываем вывод команды к нашим выводам
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func deleteCgroup(cgroupDirectory string) {
	if err := os.RemoveAll(cgroupDirectory); err != nil {
		fmt.Printf("Ошибка при очистке сgroup: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Директория cgroup успешно удалена: %s\n", cgroupDirectory)
}

func main() {
	cgroupUUID := uuid.New()
	cgroupName := "container_runner_" + cgroupUUID.String()
	cgroupDirectory := filepath.Join("/sys/fs/cgroup", cgroupName)
	memoryLimit := flag.String("memory", "64M", "Лимит по памяти")
	cpuLimit := flag.String("cpu", "50000", "Лимит по CPU в микросекундах (например, 200000 для 200 мс)")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Использование: cgroup_runner [-memory <лимит по памяти>] [-cpu <лимит по CPU>] <команда> <аргументы>")
		os.Exit(1)
	}

	command := flag.Args()[0]
	args := flag.Args()[1:]

	fmt.Printf("Название cgroup: %s\n"+
		"Полный путь до cgroup: %s\n"+
		"Введена команда: %s\n"+
		"Аргументы: %v\n"+
		"Ограничение по памяти: %s\n"+
		"Ограничение по CPU: %s\n\n",
		cgroupName, cgroupDirectory, command, args, *memoryLimit, *cpuLimit)

	if err := createCgroup(cgroupDirectory); err != nil {
		fmt.Printf("Ошибка при создании cgroup: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Cоздана cgroup \n")

	if err := setLimits(cgroupDirectory, *memoryLimit, *cpuLimit); err != nil {
		fmt.Printf("Ошибка при установлении лимитов: %v\n", err)
		deleteCgroup(cgroupDirectory)
		os.Exit(1)
	}

	fmt.Printf("Заданы лимиты cgroup\n")

	if err := bindCurrentProcessToCgroup(cgroupDirectory); err != nil {
		fmt.Printf("Не удалось привязать процесс к cgroup: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Для выхода из команды нажмите ctrl+C\n")

	if err := runCommandInNamespace(command, args); err != nil {
		fmt.Printf("Ошибка при запуске команды: %v\n", err)
		os.Exit(1)
	}

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// Ждем пока пользователь наиграется со своей командой у нас в cgroup'е
	<-exit

	deleteCgroup(cgroupDirectory)
}
