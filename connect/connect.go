package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "net"
    "os"
    "strings"
)

func main() {
    // Define server address
    serverAddr := "localhost:6666" 

    // Connect to the server
    conn, err := net.Dial("tcp", serverAddr)
    if err != nil {
        log.Fatalf("Unable to connect to server at %s: %v", serverAddr, err)
    }
    defer conn.Close()
    fmt.Printf("Connected to server at %s\n", serverAddr)

    // Create a channel to handle graceful shutdown
    done := make(chan struct{})

    // Start a goroutine to listen for messages from the server
    go func() {
        scanner := bufio.NewScanner(conn)
        for scanner.Scan() {
            message := scanner.Text()
            fmt.Println(message)
        }
        if err := scanner.Err(); err != nil && err != io.EOF {
            log.Printf("Error reading from server: %v", err)
        }
        fmt.Println("Disconnected from server.")
        close(done)
    }()

    // Read user input from the terminal
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("> ") 
        input, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                fmt.Println("\nExiting client.")
                break
            }
            log.Printf("Error reading input: %v", err)
            continue
        }

        input = strings.TrimSpace(input)
        if input == "" {
            continue 
        }

        // Send the input to the server
        _, err = fmt.Fprintln(conn, input)
        if err != nil {
            log.Printf("Error sending message: %v", err)
            break
        }
    }

    // Wait for the goroutine to finish
    <-done
}
