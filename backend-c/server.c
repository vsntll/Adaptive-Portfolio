#include <winsock2.h>
#include <ws2tcpip.h>
#include <windows.h>
#include <stdio.h>
#include <string.h>

#pragma comment(lib, "Ws2_32.lib") // Link with Ws2_32.lib

#define PORT 8080

int main() {
    WSADATA wsa;
    SOCKET server_fd, new_socket;
    struct sockaddr_in address;
    int addrlen = sizeof(address);
    char buffer[30000] = {0};
    char *header =
        "HTTP/1.1 200 OK\nContent-Type: application/json\nAccess-Control-Allow-Origin: *\n\n";

    // Initialize Winsock
    if (WSAStartup(MAKEWORD(2,2), &wsa) != 0) {
        printf("Failed to initialize Winsock. Error Code : %d\n", WSAGetLastError());
        return 1;
    }

    // Open data.json
    FILE *fp = fopen("data.json", "r");
    if (!fp) {
        perror("File open failed");
        WSACleanup();
        return 1;
    }
    char json[4096] = {0};
    fread(json, 1, sizeof(json), fp);
    fclose(fp);

    // Create socket
    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (server_fd == INVALID_SOCKET) {
        printf("Could not create socket: %d\n", WSAGetLastError());
        WSACleanup();
        return 1;
    }
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    // Bind
    if (bind(server_fd, (struct sockaddr *)&address, sizeof(address)) == SOCKET_ERROR) {
        printf("Bind failed: %d\n", WSAGetLastError());
        closesocket(server_fd);
        WSACleanup();
        return 1;
    }
    listen(server_fd, 3);

    printf("C API listening on http://localhost:%d\n", PORT);

    while (1) {
        new_socket = accept(server_fd, (struct sockaddr *)&address, &addrlen);
        if (new_socket == INVALID_SOCKET) {
            printf("Accept failed: %d\n", WSAGetLastError());
            break;
        }
        recv(new_socket, buffer, sizeof(buffer), 0);
        send(new_socket, header, strlen(header), 0);
        send(new_socket, json, strlen(json), 0);
        closesocket(new_socket);
    }
    closesocket(server_fd);
    WSACleanup();
    return 0;
}
