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
    char *hello = "HTTP/1.1 200 OK\nContent-Type: text/html\n\n";

    // Read your HTML file
    FILE *fp = fopen("index.html", "r");
    char html;
    fread(html, 1, sizeof(html), fp);
    fclose(fp);

    // Set up socket
    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    bind(server_fd, (struct sockaddr *)&address, sizeof(address));
    listen(server_fd, 3);

    while (1) {
        new_socket = accept(server_fd, (struct sockaddr *)&address, (socklen_t*)&addrlen);
        read(new_socket, buffer, 30000);
        write(new_socket, hello, strlen(hello));
        write(new_socket, html, strlen(html));
        close(new_socket);
    }
    return 0;
}
