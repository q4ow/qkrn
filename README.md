# qkrn - Distributed Key-Value Store

A simple, distributed key-value store implemented in Go with HTTP API support.

## Features

- **Thread-safe in-memory storage** with concurrent read/write support
- **HTTP REST API** for easy client integration
- **Configurable server settings** via command-line flags
- **Graceful shutdown** handling
- **Comprehensive testing** with unit tests and API tests

## Quick Start

### Prerequisites

- Go 1.24.3 or later
- `make` (optional, for using Makefile targets)

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
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
  -address string
        Node address (default "localhost")
  -log-level string
        Log level (debug, info, warn, error) (default "info")
  -node-id string
        Node ID (default hostname)
  -port int
        Node port (default 8080)
```

### API Examples

Store a key-value pair:
```bash
curl -X PUT http://localhost:8080/kv/hello \
  -H "Content-Type: application/json" \
  -d '{"value":"world"}'
```

Retrieve a value:
```bash
curl http://localhost:8080/kv/hello
```

List all keys:
```bash
curl http://localhost:8080/keys
```

Delete a key:
```bash
curl -X DELETE http://localhost:8080/kv/hello
```

## Development

### Project Structure

```
qkrn/
├── cmd/qkrn/           # Main application entry point
├── internal/           # Private application code
│   ├── api/            # HTTP API server
│   ├── config/         # Configuration management
│   └── store/          # Key-value store implementation
├── pkg/types/          # Public types and interfaces
├── docs/              # Documentation
├── scripts/           # Utility scripts
└── bin/               # Build artifacts
```

### Available Make Targets

- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run unit tests
- `make test-api` - Run API integration tests (requires server running)
- `make check` - Run all checks (format, vet, test)
- `make clean` - Clean build artifacts
- `make release` - Build release binaries for Linux and macOS

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

## Architecture

The current implementation provides a foundation for a distributed key-value store:

- **Store Layer** (`internal/store`): Thread-safe in-memory storage
- **API Layer** (`internal/api`): HTTP REST API server
- **Configuration** (`internal/config`): Command-line flag parsing
- **Types** (`pkg/types`): Shared interfaces and data structures

## Roadmap

This is the initial implementation focusing on a single-node key-value store. Future enhancements will include:

1. **Multi-node replication** - Data replication across multiple nodes
2. **Node discovery** - Automatic peer discovery and cluster formation
3. **Consensus algorithm** - Implement Raft for consistency
4. **Persistent storage** - Add disk-based storage options
5. **Advanced features** - Authentication, encryption, monitoring

See `project.md` for detailed project goals and architecture plans.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests (`make check`)
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
yet another key/value store
