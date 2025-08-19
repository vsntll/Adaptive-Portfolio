#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>

#define PORT 8080

int main() {
    int server_fd, new_socket;
    struct sockaddr_in address;
    int addrlen = sizeof(address);
    char buffer[30000] = {0};
    char *header = 
        "HTTP/1.1 200 OK\nContent-Type: application/json\nAccess-Control-Allow-Origin: *\n\n";

    FILE *fp = fopen("data.json", "r");
    if (!fp) {
        perror("File open failed");
        return 1;
    }
    char json[4096] = {0};
    fread(json, 1, sizeof(json), fp);
    fclose(fp);

    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    bind(server_fd, (struct sockaddr *)&address, sizeof(address));
    listen(server_fd, 3);

    printf("C API listening on http://localhost:%d\n", PORT);

    while (1) {
        new_socket = accept(server_fd, (struct sockaddr *)&address, (socklen_t*)&addrlen);
        read(new_socket, buffer, 30000);
        write(new_socket, header, strlen(header));
        write(new_socket, json, strlen(json));
        close(new_socket);
    }
    return 0;
}
