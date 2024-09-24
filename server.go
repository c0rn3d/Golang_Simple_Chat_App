package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "net"
    "strings"
    "sync"
    "time"
)

// Client struct holds the connection and user information for each client
type Client struct {
    conn     net.Conn // Connection object for the client
    username string   // Username of the client
    color    string   // Color for displaying messages
    lastMsg  string   // To track the last displayed message
}

// Global variables
var (
    clients     = make(map[string]Client) // Map of connected clients
    clientsLock sync.Mutex                  // Mutex for synchronizing access to the clients map
    messageChan = make(chan string)         // Channel for broadcasting messages
)

func main() {
    rand.Seed(time.Now().UnixNano()) // Seed random number generator for rolling dice

    // Start listening for TCP connections on port 8080
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer ln.Close() // Ensure the listener is closed when the function exits

    go broadcastMessages() // Start the message broadcasting goroutine

    fmt.Println("Chat server started on :8080")
    for {
        // Accept new connections
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        // Handle each new connection in a separate goroutine
        go handleConnection(conn)
    }
}

// handleConnection manages communication with a connected client
func handleConnection(conn net.Conn) {
    defer conn.Close() // Ensure the connection is closed when the function exits

    var username string
    color := "\033[0m" // Default color (reset)

    // Get username with a check for uniqueness
    for {
        fmt.Fprint(conn, "Enter your username: ") // Prompt for username
        scanner := bufio.NewScanner(conn) // Create a new scanner for reading input
        scanner.Scan() // Read the username input
        username = scanner.Text() // Get the text from the scanner

        clientsLock.Lock() // Lock the clients map for safe access
        if _, exists := clients[username]; exists { // Check if username is taken
            clientsLock.Unlock() // Unlock before continuing to prompt for a new username
            _, _ = fmt.Fprint(conn, "Username is already taken. Please choose a different username.\n")
            continue // Prompt for username again
        }
        // Register the new client
        clients[username] = Client{conn: conn, username: username, color: color}
        clientsLock.Unlock() // Unlock after modifying the clients map
        break // Exit the loop when a valid username is set
    }

    fmt.Printf("%s joined the chat\n", username) // Log that the user has joined
    messageChan <- fmt.Sprintf("%s has joined the chat\n", username) // Notify others that the user joined
    _, _ = fmt.Fprint(conn, helpMessage()) // Send help message to the new client

    // Listen for messages from this client
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        message := scanner.Text() // Read the incoming message
        if message == "/exit" { // Check for exit command
            break // Exit the loop to clean up
        }
        if message == "/clear" { // Check for clear screen command
            _, _ = fmt.Fprint(conn, "\033[H\033[2J") // Clear the terminal screen
            continue
        }
        if message == "/help" { // Check for help command
            _, _ = fmt.Fprint(conn, helpMessage()) // Send the help message
            continue
        }
        if message == "/users" { // Check for users command
            _, _ = fmt.Fprint(conn, listUsers()) // Send the list of connected users
            continue
        }
        if message == "/roll" { // Check for roll command
            roll := rand.Intn(6) + 1 // Roll a dice (1 to 6)
            messageChan <- fmt.Sprintf("%s rolled a %d\n", username, roll) // Notify others of the roll
            continue
        }
        // Handle color change command
        if strings.HasPrefix(message, "/color->") {
            newColor := message[len("/color->"):] // Extract color name from command
            switch strings.ToLower(newColor) { // Set color based on input
            case "red":
                color = "\033[31m" // Red
            case "green":
                color = "\033[32m" // Green
            case "blue":
                color = "\033[34m" // Blue
            default:
                _, _ = fmt.Fprint(conn, "Invalid color. Available options: red, green, blue.\n")
                continue
            }
            // Update client color
            clientsLock.Lock()
            client := clients[username]
            client.color = color // Update the color field
            clients[username] = client // Save back to the clients map
            clientsLock.Unlock()
            _, _ = fmt.Fprint(conn, fmt.Sprintf("Your chat color is now %s%s.\n", color, newColor))
            continue
        }

        // If the command is unrecognized, inform the user
        if strings.HasPrefix(message, "/") {
            _, _ = fmt.Fprint(conn, "Invalid command. Type /help for the list of commands.\n")
            continue
        }

        // Prepare message with color
        coloredMessage := fmt.Sprintf("%s: %s\n", username, color+message+"\033[0m")
        // Check if the last message is the same to avoid repetition
        if clients[username].lastMsg != coloredMessage {
            messageChan <- coloredMessage // Send message to the broadcast channel
            clientsLock.Lock()
            client := clients[username]
            client.lastMsg = coloredMessage // Update last displayed message
            clients[username] = client // Save back to the clients map
            clientsLock.Unlock()
        }
    }

    // Cleanup on exit
    clientsLock.Lock()
    delete(clients, username) // Remove the client from the map
    clientsLock.Unlock()
    messageChan <- fmt.Sprintf("%s has left the chat\n", username) // Notify others that the user has left
    fmt.Printf("%s disconnected\n", username) // Log disconnection
}

// broadcastMessages listens for messages on the channel and sends them to all clients
func broadcastMessages() {
    for message := range messageChan { // Loop indefinitely until the channel is closed
        clientsLock.Lock()
        // Send the message to each connected client
        for _, client := range clients {
            _, err := fmt.Fprint(client.conn, message) // Send the message
            if err != nil {
                fmt.Println("Error sending message to client:", err) // Log any errors
            }
        }
        clientsLock.Unlock() // Unlock after broadcasting the message
    }
}

// helpMessage returns a string with the available commands
func helpMessage() string {
    // Define purple color
    purple := "\033[35m" // Purple color
    reset := "\033[0m"   // Reset color
    return purple + "Available commands:\n" +
        "/help  - List available commands\n" +
        "/clear - Clear the chat screen\n" +
        "/users - List connected users\n" +
        "/roll  - Roll a dice (1-6)\n" +
        "/color->red   - Set your chat color to red\n" +
        "/color->green - Set your chat color to green\n" +
        "/color->blue  - Set your chat color to blue\n" +
        "/exit  - Disconnect from the chat\n" + reset
}

// listUsers returns a string with the current connected users
func listUsers() string {
    clientsLock.Lock() // Lock access to the clients map
    defer clientsLock.Unlock() // Ensure map is unlocked when done

    if len(clients) == 0 {
        return "No users currently connected.\n" // Return message if no users are connected
    }

    usernames := make([]string, 0, len(clients)) // Create a slice to hold usernames
    for username := range clients {
        usernames = append(usernames, username) // Add each username to the slice
    }
    return "Connected users: " + strings.Join(usernames, ", ") + "\n" // Join and return the usernames
}
