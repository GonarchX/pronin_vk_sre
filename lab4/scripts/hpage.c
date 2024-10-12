#include <stdio.h>
#include <stdlib.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>

#define SIZE (2 * 1024 * 1024) // 2 MB

int main() {
    int fd = open("/mnt/huge/my_memory_file", O_CREAT | O_RDWR, 0666);
    ftruncate(fd, SIZE);

    void *ptr = mmap(NULL, SIZE, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0);
    if (ptr == MAP_FAILED) {
        perror("mmap");
        exit(EXIT_FAILURE);
    }

    strcpy(ptr, "Hello, Huge Pages!");

    munmap(ptr, SIZE);
    close(fd);

    return 0;
}
