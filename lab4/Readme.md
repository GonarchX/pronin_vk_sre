# 1 Задание

Первым делом делаем скрипт `cpu_info.sh`, для вывода информации о CPU памяти исполняемым файлом

![chmod script.png](./imgs/1/chmod_script.png)

Затем запустим его

![cpu script.png](./imgs/1/cpu_info.png)

Информации о памяти нету, связываю это с тем, что работа велась на виртуальной машине, если вызвать отдельно команду `dmidecode --type memory`, то никакой информации не увидим

![demicode_fail.png](./imgs/1/demicode_fail.png)

Вывод должен быть похож на следующий 

![dmidecode_memory.png](./imgs/1/dmidecode_memory.png)

# 2 Задание

Первым делом замаунтим раздел под hugepages с размером 2МБ

![1.mount.png](./imgs/2/1.mount.png)

Затем скомпилируем программу

![3.compile_program.png](./imgs/2/3.compile_program.png)

Посмотрим на использование huge pages

![5.used_huge_page.png](./imgs/2/5.used_huge_page.png)

Если увеличить размер аллоцированной памяти до 8 МБ, то в выводе использования huge pages увидим 3 разервированных страниц

![6.reserved_pages.png](./imgs/2/6.reserved_pages.png)

# 3 Задание

Запустим приложение балванку, которое будет потреблять ресурсы

```
#include <stdio.h>

int main() {
    while (1);
    return 0;
}
```

Посмотрим на процент использования ресурсов при помощи `htop`

![2.top.png](./imgs/3/2.top.png)

Поменяем максимальный процент использования процессора до 75%

![3.change_cpu_usage.png](./imgs/3/3.change_cpu_usage.png)

И повысим приоритет нашего приложения 

![5.chrt.png](./imgs/3/5.chrt.png)

Запустим параллельно еще один экземпляр нашего приложения и выведем еще раз использование ресурсов при помощи `htop`

![4.top.png](./imgs/3/4.top.png)

Видим, что первому приложению отдается ~75% процессорного времени, в то время как второе приложение использует оставшиеся ресурсы.

![6.taskset.png](./imgs/3/6.taskset.png)

# 4 Задание

Создадим raid массив с `/dev/sdb1` и `/dev/sdc1`  

![1.mdadm.png](./imgs/4/1.mdadm.png)

Добавим в него `/dev/sdd`

![2.cat_proc.png](./imgs/4/2.cat_proc.png)

"Сломаем" `/dev/sdb1`

![3.fail_disk.png](./imgs/4/3.fail_disk.png)

Видим, что диск поменял статус на `Failed`, однако это никак не отразилось на нашем RAID массиве

![4.cap_proc.png](./imgs/4/4.cap_proc.png)

Перезапустим сломанный диск и увидим, что он обратно вернулся в строй

![5.repair_disk.png](./imgs/4/5.repair_disk.png)

