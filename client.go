package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
)

// Client represents a connected client
type Client struct {
    conn     net.Conn
    server   *Server
    nickname string
    // Channel allowing to queue messages to be sent to client 
    send     chan string
}

// NewClient initializes a new Client instance
func NewClient(conn net.Conn, server *Server) *Client {
    return &Client{
        conn:   conn,
        server: server,
        send:   make(chan string),
    }
}

// Handle manages the client's lifecycle.
func (c *Client) Handle() {
    defer c.conn.Close()

    // Start a goroutine to send messages to the client
    go c.writePump()

    scanner := bufio.NewScanner(c.conn)
    for scanner.Scan() {
        input := scanner.Text()
        c.processInput(input)
    }

    // Cleanup on client disconnect
    c.server.unregister <- c
    if c.nickname != "" {
        log.Printf("Client %s disconnected", c.nickname)
    } 
}

// Sends messages from the send channel to the client (Consumer)
func (c *Client) writePump() {
    // Loop over ever message in the send channel associated with each client
    for msg := range c.send {
        _, err := fmt.Fprint(c.conn, msg)
        if err != nil {
            log.Printf("Error sending message to %s: %v", c.nickname, err)
            return
        }
    }
}

// Queue a message to be sent to the client
func (c *Client) SendMessage(message string) {
    c.send <- message
}

// Pasre and handle client input
func (c *Client) processInput(input string) {
    input = strings.TrimSpace(input)
    if input == "" {
        return
    }

    if strings.HasPrefix(input, "/") {
        c.handleCommand(input)
    } else {
        if c.nickname == "" {
            c.SendMessage("Error: A command must be specified\n")
            return
        }
        // Not a command, create message 
        c.server.broadcast <- Message{
            From:    c.nickname,
            Content: input,
        }
    }
}

// Execute client commands 
func (c *Client) handleCommand(input string) {
    // Splits input at white space, assigns each part to an element in parts array
    parts := strings.SplitN(input, " ", 2)
    cmd := strings.ToUpper(parts[0])
    args := ""
    if len(parts) > 1 {
        args = parts[1]
    }

    switch cmd {
    case "/NICK", "/N":
        c.setNickname(args)
    case "/SEND", "/S":
        c.sendDirectMessage(args)
    case "/BCAST", "/B":
        c.broadcastMessage(args)
    case "/LIST", "/L":
        c.listUsers()
    default:
        c.SendMessage("ERROR: Invalid command. Valid commands are:\n" +
            "/NICK or /N <name>\n" +
            "/LIST or /L\n" +
            "/BCAST or /B <message>\n" +
            "/SEND or /S <nickname(s)> <message>\n")
    }
}

// Sets the client's nickname after validation
func (c *Client) setNickname(args string) {
    args = strings.TrimSpace(args)
    if args == "" {
        c.SendMessage("Usage: /NICK <name>\n")
        return
    }

    // Only the first argument is considered
    nick := strings.Fields(args)[0] 
    if !isValidNickname(nick) {
        c.SendMessage("Invalid nickname. Must start with a letter and contain only letters, numbers, or underscores. Max 10 characters.\n")
        return
    }

    // Lock since accessing shared client map
    c.server.clientsLock.Lock()
    defer c.server.clientsLock.Unlock()

    if _, exists := c.server.clients[nick]; exists {
        c.SendMessage("Error: Nickname is already in use!\n")
        return
    }

    // Change nickname if client already has an existing one
    if c.nickname != "" {
        delete(c.server.clients, c.nickname)
    }

    c.nickname = nick
    // Add client to clients map
    c.server.clients[nick] = c
    c.server.register <- c
    c.SendMessage(fmt.Sprintf("Nickname successfully set to %s\n", nick))
}

// Sends a message to specified nicknames
func (c *Client) sendDirectMessage(args string) {
    if c.nickname == "" {
        c.SendMessage("You must set a nickname before using /SEND.\n")
        return
    }
    
    // Break args into fields based on white space
    fields := strings.Fields(args)
    if len(fields) < 2 {
        c.SendMessage("Invalid format. Usage: /SEND <nickname(s)> <message>\n")
        return
    }

    recipientsPart := fields[0]
    message := strings.Join(fields[1:], " ")

    recipients := strings.Split(recipientsPart, ";")

    // Send message to each recipient
    for _, recipient := range recipients {
        recipient = strings.TrimSpace(recipient)
        if recipient == "" {
            continue
        }

        c.server.broadcast <- Message{
            From:    c.nickname,
            To:      recipient,
            Content: message,
        }

        c.SendMessage(fmt.Sprintf("Message sent to %s\n", recipient))
    }
}

// Sends a message to all connected clients.
func (c *Client) broadcastMessage(args string) {
    if c.nickname == "" {
        c.SendMessage("You must set a nickname before using /BCAST.\n")
        return
    }

    message := strings.TrimSpace(args)
    if message == "" {
        c.SendMessage("Usage: /BCAST <message>\n")
        return
    }

    c.server.broadcast <- Message{
        From:    c.nickname,
        Content: message,
    }

    c.SendMessage("Broadcast message sent.\n")
}

// Sends the list of connected users to the client.
func (c *Client) listUsers() {
    // Lock since we are accessing client map 
    c.server.clientsLock.RLock()
    defer c.server.clientsLock.RUnlock()

    if len(c.server.clients) == 0 {
        c.SendMessage("No users currently connected.\n")
        return
    }

    var users []string
    for nickname := range c.server.clients {
        users = append(users, nickname)
    }

    c.SendMessage("Connected users: " + strings.Join(users, ", ") + "\n")
}
