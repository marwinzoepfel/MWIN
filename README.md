# MWIN
**MWIN** (My Witty Interactive Nonsense): A simple, self-hosted, terminal-based chat room application written in Go. No installation required, just download and run!

a Simple Go Chat Room

This is a basic, terminal-based chat room implemented in Go. It allows multiple users to connect and chat with each other in real time.

## Features

- **Multi-User Chat:** Multiple users can join and participate in the chat room simultaneously.
- **Real-Time Messaging:** Messages are instantly broadcast to all connected users.
- **Customizable Usernames:** Users can choose their own name or use the default (hostname).
- **Timestamps:** Messages are displayed with timestamps for better clarity.
- **Customizable Port:** The server can be configured to run on a port of your choice.

## Prerequisites

To run this chat application, you don't need any prerequisites. Simply download the appropriate binary for your operating system and follow the usage instructions below.

## Prerequisites for Developers

If you want to modify or customize the chat application, you'll need Go installed on your system. You can download it from the official website: [https://golang.org/](https://golang.org/)

## Usage

1. **Server:**
   - Download `chat_server` for your operating system.
   - Open a terminal and navigate to the directory containing `chat_server`.
   - Run `./chat_server` (or `chat_server.exe` on Windows).
   - You'll be prompted to enter the desired port (default: 8080). If you press Enter without entering a port, the default will be used.

2. **Client:**
   - Download `chat_client` for your operating system.
   - Open a terminal and navigate to the directory containing `chat_client`.
   - Run `./chat_client` (or `chat_client.exe` on Windows).
   - You'll be prompted to enter the server address (e.g., `localhost:8080`) and your desired name (default: your device's hostname).

3. **Chat:**
   - Once connected, type your messages in the client terminal and press Enter to send them.
   - Messages from all connected users will be displayed in each client's terminal, along with timestamps and usernames.

## Configuration

- **Default Port:** The default port (8080) is defined in the `defaultPort` constant in `server/main.go`. You can modify this value if needed.

## Building from Source (For Developers)

1. **Server:** Navigate to the `server` directory and run `go build` to create the `chat_server` executable.
2. **Client:** Navigate to the `client` directory and run `go build` to create the `chat_client` executable.

## Contributing

Feel free to fork this repository and submit pull requests to add new features, fix bugs, or improve the code.

## License

This project is licensed under the GNU Affero General Public License version 3 (AGPLv3) - see the LICENSE file for details.

**Ethical Considerations**

By using this software, you agree to uphold the following ethical principles:

- **Do No Harm:** This software should not be used to cause harm, exploit vulnerabilities, or engage in any malicious activities.
- **Respect Privacy:** If this software collects or processes data, it should be done in a transparent and respectful manner, prioritizing user privacy and consent.
- **Promote Inclusivity:** This software should be accessible and inclusive, avoiding any discriminatory or exclusionary practices.
- **Open Collaboration:**  We encourage collaboration and contributions to this project that align with these ethical guidelines.

**Contributions**

We welcome contributions that enhance the functionality and ethical integrity of this project. Please review our CONTRIBUTING.md guidelines before submitting any changes.

**Disclaimer**

This software is provided "as is", without warranty of any kind, express or implied, including but not limited to the warranties of merchantability, fitness for a particular purpose and noninfringement. In no event shall the authors or copyright holders be liable for any claim, damages or other liability, whether in an action of contract, tort or otherwise, arising from, out of or in connection with the software or the use or other dealings in the software.
