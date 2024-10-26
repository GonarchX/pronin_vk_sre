import os
import sys
import signal
import subprocess
import argparse
import tempfile
import uuid

MEMORY_LIMIT_FILE = "memory.max"
CPU_LIMIT_FILE = "cpu.max"
PROCS_FILE = "cgroup.procs"

def create_cgroup(cgroup_directory):
    try:
        os.makedirs(cgroup_directory, exist_ok=True)
    except Exception as e:
        raise Exception(f"Ошибка при создании директории {cgroup_directory}: {e}")

def set_limits(cgroup_directory, memory_limit, cpu_limit):
    try:
        with open(os.path.join(cgroup_directory, MEMORY_LIMIT_FILE), 'w') as f:
            f.write(memory_limit)
        with open(os.path.join(cgroup_directory, CPU_LIMIT_FILE), 'w') as f:
            f.write(cpu_limit)
    except Exception as e:
        raise Exception(f"Ошибка при задании лимитов: {e}")

def bind_current_process_to_cgroup(cgroup_directory):
    try:
        with open(os.path.join(cgroup_directory, PROCS_FILE), 'w') as f:
            f.write(str(os.getpid()))
    except Exception as e:
        raise Exception(f"Ошибка при привязке процесса к cgroup: {e}")

def run_command_in_namespace(command, args, root_directory):
    subprocess.run(["unshare", "--pid", "--fork", "--mount-proc", "--root", root_directory, command] + args)

def delete_cgroup(cgroup_directory):
    try:
        os.rmdir(cgroup_directory)
        print(f"Директория cgroup успешно удалена: {cgroup_directory}")
    except Exception as e:
        print(f"Ошибка при очистке cgroup: {e}")

def main():
    parser = argparse.ArgumentParser(description='Запуск команды в cgroup с ограничениями по ресурсам')
    parser.add_argument('-memory', default='64M', help='Лимит по памяти')
    parser.add_argument('-cpu', default='50000', help='Лимит по CPU в микросекундах')
    parser.add_argument('-root', default='/mnt/rootfs', help='Корневая директория для пространства имен')
    parser.add_argument('command', help='Команда для выполнения')
    parser.add_argument('args', nargs=argparse.REMAINDER, help='Аргументы команды')

    args = parser.parse_args()

    cgroup_uuid = uuid.uuid4()
    cgroup_name = f'container_runner_{cgroup_uuid}'
    cgroup_directory = os.path.join('/sys/fs/cgroup', cgroup_name)

    print(f"Название cgroup: {cgroup_name}")
    print(f"Полный путь до cgroup: {cgroup_directory}")
    print(f"Введена команда: {args.command}")
    print(f"Аргументы: {args.args}")
    print(f"Ограничение по памяти: {args.memory}")
    print(f"Ограничение по CPU: {args.cpu}")
    print(f"Корневая директория: {args.root}\n")

    create_cgroup(cgroup_directory)
    print("Создана cgroup")

    set_limits(cgroup_directory, args.memory, args.cpu)
    print("Заданы лимиты cgroup")

    bind_current_process_to_cgroup(cgroup_directory)
    print("Процесс привязан к cgroup")

    try:
        run_command_in_namespace(args.command, args.args, args.root)
    except Exception as e:
        print(f"Ошибка при запуске команды: {e}")
        sys.exit(1)

if name == "main":
    main()