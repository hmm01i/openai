# OpenAI CLI Client

A powerful command-line interface for interacting with OpenAI's services, including chat models like GPT-4 and image generation capabilities.

## Features

- 💬 Interactive chat with OpenAI models (GPT-4, GPT-3.5, etc.)
- 🎭 Persona management for different AI roles
- 🖼️ Image generation from text descriptions
- 💾 Conversation history management
- 🌐 HTTP server mode for API access
- 🔐 Secure API token handling

## Installation

### Prerequisites

- Go 1.20 or higher
- OpenAI API key

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/hmm01i/openai.git
cd openai
```

2. Build the binary:
```bash
make build
```

3. (Optional) Install system-wide:
```bash
make install
```

## Configuration

1. Set your OpenAI API token in one of two ways:
   - Environment variable: `export OPENAI_API_TOKEN=your_token_here`
   - Token file: Create `~/.openai/token` and paste your token there

The application will create the following directory structure:
```
~/.openai/
├── token           # API token file
├── personas/       # Saved AI personas
└── conversations/  # Saved conversations
```

## Usage

### Basic Commands

Start the chat interface:
```bash
oai chat
```

Start the HTTP server:
```bash
oai chat server
```

Generate an image:
```bash
oai image -p "your image description" -o output.png
```

### Chat Commands

While in chat mode, you can use these commands:

- `/help` - Show available commands
- `/persona` - Manage AI personas
  - `list` - List all personas
  - `show` - Show current persona
  - `save <name>` - Save current system directive as a persona
  - `load <name>` - Load a persona
- `/model` - Manage AI models
  - `list` - List available models
  - `set <model>` - Set current model
- `/history` - Manage chat history
  - `show` - Display conversation history
  - `clear` - Clear current history
- `/conversation` - Manage conversations
  - `list` - List saved conversations
  - `save <name>` - Save current conversation
  - `load <name>` - Load a conversation
- `/system` - System commands
  - `directive <text>` - Set system directive
- `/q` - Quit the application

### HTTP Server Mode

When running in server mode, the following endpoints are available:

- `POST /chat` - Send chat messages
- `POST /syscmd` - Execute system commands

The server runs on port 8080 by default.

## Development

### Project Structure

```
.
├── cmd/
│   ├── main.go       # Main entry point
│   ├── chat.go       # Chat functionality
│   ├── image.go      # Image generation
│   └── api.go        # HTTP server
├── pkg/
│   ├── commands/     # Command system
│   └── version/      # Version information
└── Makefile         # Build configuration
```

### Make Commands

- `make build` - Build the binary
- `make clean` - Clean build artifacts
- `make test` - Run tests
- `make install` - Install the binary
- `make tidy` - Run go mod tidy

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.