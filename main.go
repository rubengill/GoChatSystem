package main

import (
    "log"
    "net"
)

func main() {
    server := NewServer()
    go server.Run()

    port := "6666"
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        log.Fatalf("Failed to listen on port %s: %v", port, err)
    }
    defer listener.Close()
    log.Printf("Server is listening on port %s", port)

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }
        log.Printf("Accepted connection from %s", conn.RemoteAddr())
        client := NewClient(conn, server)
        go client.Handle()
    }
}
