package main

import (
    "log"
    "sync"
    "fmt"
)

// Server manages all connected clients and their nicknames.
type Server struct {
    clients     map[string]*Client // nickname -> client
    clientsLock sync.RWMutex
    register    chan *Client
    unregister  chan *Client
    broadcast   chan Message
}

// Message represents a chat message.
type Message struct {
    From    string
    To      string
    Content string
}

// NewServer initializes a new Server instance.
func NewServer() *Server {
    return &Server{
        clients:    make(map[string]*Client),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan Message),
    }
}

// Run starts the server's main loop.
func (s *Server) Run() {
    for {
        // Handle whichever channel recieves data first
        select {
        case client := <-s.register:
            s.clientsLock.Lock()
            s.clients[client.nickname] = client
            s.clientsLock.Unlock()
            log.Printf("Client %s registered with nickname '%s'", client.conn.RemoteAddr(), client.nickname)
        case client := <-s.unregister:
            s.clientsLock.Lock()
            if client.nickname != "" {
                delete(s.clients, client.nickname)
                log.Printf("Client %s with nickname '%s' unregistered", client.conn.RemoteAddr(), client.nickname)
            }
            s.clientsLock.Unlock()
        case msg := <-s.broadcast:
            if msg.To == "" {
                // Broadcast message
                s.clientsLock.RLock()
                for nickname, client := range s.clients {
                    if nickname != msg.From {
                        client.SendMessage(formatMessage(msg.From, msg.Content))
                    }
                }
                s.clientsLock.RUnlock()
            } else {
                // Direct message
                s.clientsLock.RLock()
                if recipient, ok := s.clients[msg.To]; ok {
                    recipient.SendMessage(formatMessage(msg.From, msg.Content + "\n"))
                }  else {
                    // Inform the sender that the recipient is not registered
                    if sender, ok := s.clients[msg.From]; ok {
                        sender.SendMessage(fmt.Sprintf("Error: %s is not registered.\n", msg.To))
                    }
                }
                s.clientsLock.RUnlock()
            }
        }
    }
}

// formatMessage formats the message to be sent to clients.
func formatMessage(from, content string) string {
    return "[" + from + "]: " + content + "\n"
}
