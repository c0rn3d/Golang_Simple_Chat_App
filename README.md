Chat Server

This is a simple chat server implemented in Go, allowing users to connect and communicate using either a custom Python client or Telnet.
Features

    Supports multiple users.
    Customizable chat colors.
    Commands for rolling dice, clearing the chat, and listing connected users.
    Supports graceful disconnection.

Getting Started
Prerequisites

    Go: Required to run the server.
    Python: Required to run the client script.
    Telnet: A Telnet client (usually pre-installed on most systems).

Running the Server

    Clone this repository:

    bash

git clone <repository-url>

Navigate to the directory:


cd <directory-name>

Run the server using:


    go run server.go

Connecting to the Server

You can connect to the server using either:
1. Telnet:

Open your terminal and run:


telnet localhost 8080

2. Python Client:

Run the client.py script:

python client.py

User Instructions

    Enter Username: When prompted, enter a unique username.
    Available Commands:
        /help: List available commands.
        /clear: Clears the chat screen.
            Note:
                On Unix/Linux systems, this will clear the screen.
                On Windows, use the cls command directly in the terminal.
        /users: List connected users.
        /roll: Roll a dice (1-6).
        /color->red: Set your chat color to red.
        /color->green: Set your chat color to green.
        /color->blue: Set your chat color to blue.
        /exit: Disconnect from the chat.

Example Usage

After entering your username, you can type messages and use the available commands. For example:



Enter your username: alice
alice: Hello everyone!
/roll
alice rolled a 4
/users
Connected users: alice, bob
/clear

Disconnection

To disconnect from the chat, use the /exit command.
Contributing

Feel free to contribute to this project by submitting a pull request. Suggestions and improvements are welcome!
License

This project is licensed under the MIT License. See the LICENSE file for details.
