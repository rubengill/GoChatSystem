# GoChatSystem

An Go based real time messaging platform created with Go to facilitate communication between users over TCP connections. This system supports setting nicknames, listing active users, sending direct messages, and broadcasting messages to all users. A Java and Go client is provided to interact with the chat server.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [Starting the Chat Server](#starting-the-chat-server)
  - [Running the Go / Java Client](#running-the-java-client)
- [Commands](#commands)
  - [Setting a Nickname](#setting-a-nickname)
  - [Listing Active Users](#listing-active-users)
  - [Sending Direct Messages](#sending-direct-messages)
  - [Broadcasting Messages](#broadcasting-messages)
- [Testing](#testing)
- [License](#license)

## Features

- **Nickname Management**: Users can set and change their unique nicknames.
- **User Listing**: Retrieve a list of all active nicknames.
- **Direct Messaging**: Send messages to specific users.
- **Broadcast Messaging**: Send messages to all registered users.
- **Concurrent Connections**: Supports multiple simultaneous TCP connections.
- **Error Handling**: Validates commands and provides appropriate feedback.

## Architecture

The system consists of several Go components and 2 clients (You choose which to run):

1. **Server**
   - **Type**: Go service with concurrent message handling
   - **Role**:  Manages connected clients, maintains nickname mappings, and handles message routing through Go channels for broadcasting and direct messaging.
2. **Client**
   - **Type**: Connection handler
   - **Role**:  Intermediary that reads from the TCP connection, processes commands, and forwards messages to the server via channels, delivering messages to the connected user
3. **Go / Java Client**
   -  **Type**: End User terminal application
   - **Functionality**: Connects to the Go chat server over TCP, provides a text based interface for sending commands and messages, and displays incoming messages from other users in real-time.

The system leverages Go's concurrency model with goroutines and channels to efficiently handle multiple simultaneous connections and message distribution.

## Requirements

Ensure the following software is installed on your system to successfully set up and run the Chat System:

- **Go**: Version `1.21.1`
- **Java**: JDK `17`

### Installation Links

- **Go**: [Installation Guide](https://go.dev/doc/install)

> **Note:** These are the versions used during development and testing. While newer versions may be compatible, using these specific versions ensures compatibility and reduces the likelihood of encountering unexpected issues.

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/rubengill/GoChatSystem
cd GoChatSystem
```

### 2. Build the Go Project
```bash
go run *.go
```

### 3. (If using Go) Connect the client
```bash
cd connect
go run connect.go
```

### 4. (If using Java) Compile the Java Client 
> **Note:** Ensure in the root mix directory
```bash
cd connect
javac ChatClient.java
```

### 5. Run the Java Client 
> **Note:** Ensure in the root mix directory
```bash
java ChatClient
```

## Commands

Users interact with the chat system using specific commands. Commands are case-sensitive.

### Setting a Nickname

- **Command Variants**: `/NICK`, `/N`
- **Usage**: `/NICK <nickname>`
  - `<nickname>`: The nickname to register the process.

**Rules**:
- Must start with an alphabet.
- Can contain alphanumeric characters and underscores.
- Maximum length of 10 characters.
- Nicknames must be unique.
- Required before sending or receiving messages.

**Examples**:
- `/NICK homer`
- `/N homer`

**Responses**:
- **Success**: Confirmation message.
- **Failure**: Error indicating nickname is invalid or already in use.

### Listing Active Users

- **Command Variants**: `/LIST`, `/L`
- **Usage**: `/LIST`

**Description**: Retrieves a list of all currently registered nicknames.

**Examples**:
- `/LIST`
- `/L`

**Responses**:
- **Success**: Retreives all active users.
- **Failure**: Error indicating failure to retreive users.

### Sending Direct Messages

- **Command Variants**: `/SEND`, `/S`
- **Usage**: `/SEND <nicknames> <message>`
  - `<nicknames>`: One or more nicknames separated by semicolons (;).
  - `<message>`: The message to be sent.

**Examples**:
- `/SEND homer hello world`
- `/S homer;bart hello everyone`

**Responses**:
- **Success**: Message delivered confirmation.
- **Failure**: Error indicating invalid recipients.

### Broadcasting Messages

- **Command Variants**: `/BCAST`, `/B`
- **Usage**: `/BCAST <message>`

**Examples**:
- `/BCAST hello all`
- `/B this is a broadcast message`

**Responses**:
- **Success**: Message broadcast confirmation.
- **Failure**: Error sending message.

---

## Testing

To ensure the chat system functions correctly, follow these testing steps:

1. **Start the Chat Server and Proxy Server**
   - Run the following command:
     ```
     go run *.go
     ```

2. **Run Multiple Clients**
   - Open multiple terminal windows and run the Java client in each:
    > **Option A**: Go Client
    ```bash
    cd connect
    go run connect.go
    ```


   > **Option B**: Java Client
    **Note:** Ensure ChatClient.java is compiled
     ```
     java ChatClient
     ```

3. **Set Nicknames**
   - In each client, set a unique nickname:
     ```
     /NICK user1
     /NICK user2
     ```

4. **List Active Users**
   - Use the `/LIST` command to verify active nicknames:
     ```
     /LIST
     ```

5. **Send Direct Messages**
   - From `user1`, send a message to `user2`:
     ```
     /SEND user2 Hello, user2!
     ```

6. **Broadcast Messages**
   - From `user1`, broadcast a message to all users:
     ```
     /BCAST Hello!
     ```

7. **Change Nicknames**
   - Change a user's nickname and ensure the update reflects across all clients:
     ```
     /NICK newUser1
     ```

8. **Invalid Commands**
   - Test invalid commands to ensure proper error handling:
     ```
     /RANDOM
     /SEND hello
     /NICK
     ```
---

## License

This project is licensed under the MIT License.
