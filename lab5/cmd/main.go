package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
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

func runCommandInNamespace(command string, args []string, rootfs string) error {
	cmdArgs := make([]string, len(args)+1)
	cmdArgs = append(cmdArgs, []string{"--fork", "--pid" /*"--mount", "--root"*/}...)
	//cmdArgs = append(cmdArgs, rootfs)
	cmdArgs = append(cmdArgs, command)
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("unshare", cmdArgs...)

	// Привязываем вывод команды к нашим выводам
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
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
	rootfs := flag.String("rootfs", ".", "Путь к корневой файловой системе")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Использование: cgroup_runner [-memory <лимит по памяти>] [-cpu <лимит по CPU>] [-rootfs <путь к rootfs>] <команда> <аргументы>")
		os.Exit(1)
	}

	command := flag.Args()[0]
	args := flag.Args()[1:]

	fmt.Printf("Название cgroup: %s\n"+
		"Полный путь до cgroup: %s\n"+
		"Введена команда: %s\n"+
		"Аргументы: %v\n"+
		"Ограничение по памяти: %s\n"+
		"Ограничение по CPU: %s\n"+
		"Путь к rootfs: %s\n\n",
		cgroupName, cgroupDirectory, command, args, *memoryLimit, *cpuLimit, *rootfs)

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

	fmt.Printf("Для выхода из команды нажмите ctrl+C\n")

	if err := runCommandInNamespace(command, args, *rootfs); err != nil {
		fmt.Printf("Ошибка при запуске команды: %v\n", err)
		os.Exit(1)
	}

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// Ждем пока пользователь наиграется со своей командой у нас в cgroup'е
	<-exit

	deleteCgroup(cgroupDirectory)
}
