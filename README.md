# qkrn - Distributed Key-Value Store

A simple, ephemeral key-value store implemented in Go designed with REST in mind

## Features

- **Thread-safe in-memory storage** with concurrent read/write support
- **HTTP REST API** for easy client integration
- **API Key Authentication** with secure token generation and validation
- **Configurable server settings** via command-line flags
- **Gracefully handles** interruptions and poweroff signals

## Quick Start

### Prerequisites

- Go (any stable version should work)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/q4ow/qkrn
   cd qkrn
   ```

2. Build the application:
   ```bash
   make build
   # or
   go build -o bin/qkrn ./cmd/qkrn
   ```

3. Run the server:
   ```bash
   ./bin/qkrn
   # or
   make run
   ```

### Usage

The server starts on `localhost:8080` by default. You can configure it using command-line flags:

```bash
./bin/qkrn --help
Usage of ./bin/qkrn:
    ...
```

### Authentication

qkrn supports optional API key authentication:

```bash
# Start with authentication enabled and auto-generate API key
./bin/qkrn --auth-enabled

# Start with specific API key
./bin/qkrn --auth-enabled --api-key "your-secure-api-key"

# Start without authentication (default)
./bin/qkrn
```

### API Examples

Store a key-value pair:
```bash
# Without authentication
curl -X PUT http://localhost:8080/kv/hello \
  -H "Content-Type: application/json" \
  -d '{"value":"world"}'

# With authentication
curl -X PUT http://localhost:8080/kv/hello \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"value":"world"}'
```

Retrieve a value:
```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  http://localhost:8080/kv/hello
```

List all keys:
```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  http://localhost:8080/keys
```

Delete a key:
```bash
curl -X DELETE \
  -H "Authorization: Bearer YOUR_API_KEY" \
  http://localhost:8080/kv/hello
```

## Development

### Project Structure

```
qkrn/
├── cmd/qkrn/           # Main application entry point
├── internal/           # Private application code
│   ├── api/            # HTTP API server
│   ├── auth/           # Authentication middleware and utilities
│   ├── config/         # Configuration management
│   └── store/          # Key-value store implementation
├── pkg/types/          # Public types and interfaces
├── docs/              # Documentation
├── scripts/           # Utility scripts
└── bin/               # Build artifacts
```

### Running Tests

Unit tests:
```bash
make test
```

API integration tests:
```bash
# Terminal 1: Start the server
make run

# Terminal 2: Run API tests
make test-api
```

### Code Quality

The project follows Go best practices:

- Code is formatted with `go fmt`
- Code is vetted with `go vet`
- Comprehensive unit tests
- Thread-safe concurrent operations
- Proper error handling

CI can only do so much, please be nice in your pull requests :)

## Roadmap

This is just the beginning, I have much more planned for the future of qkrn:

1. **Multi-node replication** - Data replication across multiple nodes
2. **Node discovery** - Automatic peer discovery and cluster formation
3. **Consensus algorithm** - Implement Raft for consistency
4. **Persistent storage** - Add disk-based storage options
5. **Advanced features** - Authentication, encryption, monitoring

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and other make targets
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
yet another key/value store
