# Gobot AI Chat

A modern chat application that uses Ollama to provide AI-powered responses with a clean, user-friendly interface.

## Features

- Real-time AI chat interface
- Formatted responses with proper list handling
- Clean, modern UI with message bubbles
- Docker support for easy deployment
- Context-aware conversations

## Prerequisites

1. [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
2. [Ollama](https://ollama.ai) installed and running locally
3. The Mistral model pulled in Ollama

## Quick Start with Docker

1. Install Ollama from [https://ollama.ai](https://ollama.ai)

2. Pull the Mistral model:
   ```bash
   ollama pull mistral
   ```

3. Start Ollama:
   ```bash
   ollama serve
   ```

4. Build and run the application:
   ```bash
   docker-compose up --build
   ```

5. Open your browser and visit:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## Development Setup (without Docker)

If you prefer to run the application without Docker:

### Backend

1. Install Go 1.21 or later
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run the backend:
   ```bash
   go run main.go
   ```

### Frontend

1. Install Node.js and npm
2. Install dependencies:
   ```bash
   cd frontend
   npm install
   ```
3. Start the development server:
   ```bash
   npm start
   ```

## Project Structure

```
.
├── agent/         # AI agent implementation
├── frontend/      # React frontend application
├── browser/       # Web browser automation (optional)
├── input/         # Input handling utilities
├── parser/        # Query parsing utilities
├── Dockerfile     # Backend Docker configuration
├── docker-compose.yml  # Docker services orchestration
└── README.md
```

## Environment Variables

- `OLLAMA_HOST`: Host where Ollama is running (default: localhost)
- Additional environment variables can be set in `.env` files

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 