# Welcome!

This repository is a Redis clone built in Go. Working on this project helped me understand TCP server implementation, the Redis protocol, and distributed key-value servers.

## Instructions to Run

1. **Clone this repository**
   ```sh
   git clone https://github.com/arjun-gopikrishnan/redis-golang.git
   ```
2. **Install Telnet**
   - On Debian-based systems (e.g., Ubuntu):
     ```sh
     sudo apt-get install telnet
     ```
   - On Red Hat-based systems (e.g., Fedora):
     ```sh
     sudo yum install telnet
     ```
   - On macOS:
     ```sh
     brew install telnet
     ```
3. **Start the Server**
   ```sh
   go run .
   ```
4. **Connect via Telnet**
   ```sh
   telnet localhost <PORT>
   ```
   Replace `<PORT>` with the port number your server is running on.

Feel free to explore the code, suggest improvements, or use it as a learning resource. Enjoy!
