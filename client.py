import socket
import threading
import sys
import select

def receive_messages(sock):
    """Function to receive messages from the server."""
    while True:
        try:
            # Receive message from the server
            message = sock.recv(1024).decode('utf-8')
            if message:
                # Print message if it is not the username prompt
                if "Enter your username:" not in message:
                    print(message)
            else:
                break  # Exit the loop if the connection is closed
        except Exception as e:
            print("Error receiving message:", e)
            break  # Exit the loop on error

def main():
    """Main function to run the chat client."""
    # Server details
    host = 'localhost'  # Change this to your server's IP address if needed
    port = 8080         # Server port

    # Prompt for username
    username = input("Enter your username and type command /help for more info: ")

    # Create a socket object
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Connect to the server
    sock.connect((host, port))

    # Send the username to the server
    sock.sendall(f"{username}\n".encode('utf-8'))

    # Start a thread to handle incoming messages
    thread = threading.Thread(target=receive_messages, args=(sock,))
    thread.start()

    # Main loop for sending messages
    try:
        while True:
            # Use select to wait for input
            input_ready, _, _ = select.select([sys.stdin], [], [])
            for s in input_ready:
                if s == sys.stdin:
                    message = sys.stdin.readline().rstrip()  # Read input and strip newline
                    if message:
                        # Handle the exit command
                        if message == "/exit":
                            sock.sendall(f"{message}\n".encode('utf-8'))  # Inform the server
                            print("Disconnected.")
                            sock.close()  # Close the socket
                            sys.exit()  # Exit the program
                        else:
                            # Send the message to the server
                            sock.sendall(f"{message}\n".encode('utf-8'))
    except KeyboardInterrupt:
        print("Disconnected.")
        sock.close()  # Close the socket on exit
        sys.exit()  # Exit the program

if __name__ == "__main__":
    main()
