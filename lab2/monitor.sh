#!/bin/bash

CPU_THRESHOLD=80 # Значение в процентах
MEM_THRESHOLD=80 # Значение в процентах
DISK_THRESHOLD=80 # Значение в процентах
PACKET_LOSS_THRESHOLD=5 # Значение в процентах

PORT_THRESHOLD=100 # Абсолютное значение
CONNECTION_THRESHOLD=200 # Абсолютное значение

monitor() {
    while true; do
        # Текущую нагрузку на CPU
        CPU_LOAD=$(top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1}')
        
        # Текущее использование оперативной памяти
        MEM_USAGE=$(free | grep Mem | awk '{print $3/$2 * 100.0}')
        
        # Текущее использование диска
        DISK_USAGE=$(df / | grep / | awk '{ print $5 }' | sed 's/%//g')

        # Количество занятых портов
        ZON_PORT_COUNT=$(netstat -tuln | grep LISTEN | wc -l)

        # Количество активных соединений
        ACTIVE_CONNECTIONS=$(netstat -ant | grep ESTABLISHED | wc -l)

        if (( $(echo "$CPU_LOAD > $CPU_THRESHOLD" | bc -l) )); then
            echo "[WARNING] CPU usage is above threshold: ${CPU_LOAD}%"
        else
            echo "[INFO] CPU usage: ${CPU_LOAD}%"
        fi

        if (( $(echo "$MEM_USAGE > $MEM_THRESHOLD" | bc -l) )); then
            echo "[WARNING] Memory usage is above threshold (${MEM_TRESHOLD}%): ${MEM_USAGE}%"
        else
            echo "[INFO] Memory usage: ${MEM_USAGE}%"
        fi

        if (( $(echo "$DISK_USAGE > $DISK_THRESHOLD" | bc -l) )); then
            echo "[WARNING] Disk usage is above threshold (${DISK_THRESHOLD}%): ${DISK_USAGE}%"
        else
            echo "[INFO] Disk usage: ${DISK_USAGE}%"
        fi

        if (( ZON_PORT_COUNT > PORT_THRESHOLD )); then
            echo "[WARNING] Number of occupied ports is above threshold (${PORT_THRESHOLD}%): ${ZON_PORT_COUNT}"
        else
            echo "[INFO] Number of occupied ports: ${ZON_PORT_COUNT}"
        fi

        if (( ACTIVE_CONNECTIONS > CONNECTION_THRESHOLD )); then
            echo "[WARNING] Number of active connections is above threshold (${CONNECTION_THRESHOLD}%): ${ACTIVE_CONNECTIONS}"
        else 
            echo "[INFO] Current number of active connections: ${ACTIVE_CONNECTIONS}"
        fi

        # Вычисление процента потерь пакетов
        
        # Берем значения из netstat, скипаем первые две строки, т.к. в них находятся заголовки, записываем сумму полученных и потерянных в переменную total для каждого интерфейса
        netstat -i | awk 'NR>2 {total += ($3 + $4); lost += $4} END {if (total > 0) print (lost / total) * 100; else print 0}' | while read -r LOST_PERCENT; do
            if (( $(echo "$LOST_PERCENT > $PACKET_LOSS_THRESHOLD" | bc -l) )); then
                echo "[WARNING] Packet loss is above threshold (${PACKET_LOSS_THRESHOLD}%): ${LOST_PERCENT}%"
            else
                echo "[INFO] Packet loss: ${LOST_PERCENT}%"
            fi
        done
    
        PACKETS_INFO=$(netstat -i | awk 'NR>2 {total += $3; lost += $4} END {print total, lost}')
        TOTAL_PACKETS=$(echo $PACKETS_INFO | awk '{print $1}')
        LOST_PACKETS=$(echo $PACKETS_INFO | awk '{print $2}')

        # Вывод общего количества полученных и потерянных пакетов
        echo "[INFO] Total packets received: ${TOTAL_PACKETS}, Lost packets: ${LOST_PACKETS}"

        echo "---------------------------------------------------------------------- "
        
        # Чтобы не спамить логами, делаем delay в 1 секунду
        sleep 1
    done
}

monitor