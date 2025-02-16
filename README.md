# Averse Cache Server

Provide a gRCP interface for the [Goutte](https://github.com/shellkah/goutte) cache. Can be used as a standalone server or deployed as a container.

## Features

- **High Concurrency Support**: Safe for concurrent access under heavy loads.
- **LRU Eviction Policy**: Automatically removes the least recently used entries when the cache exceeds its capacity.
- **Optional TTL**: Automatically removes expired items with precision using a min-heap (priority queue) to track expiration times.
- **Fast Lookups**: Uses a hash map for O(1) average-time complexity for queries.
- **Configurable**: Easily configure the server via a config file or environment variables.
- **gRCP API**: Provides basic operations such as `Get`, `Set`, and `Delete`.

## Incoming

- **Go Client library**: gRCP wrapper for python applications.
- **Python Client library**: gRCP wrapper for Go applications.
- **Kubernetes manifest**: Simplifies deployment in a Kubernetes cluster.
- **Metrics**: Tracks cache usage and exposes it for monitoring purposes.

## Configuration
- *CACHE_CAPACITY*
- *SERVER_HOST*
- *SERVER_PORT*
- *LOG_LEVEL*

## Installation

To build and run the server from source, execute:

```bash
go mod download
go build -o ./averse cmd/server/server.go
./averse
```

## Contributing

Contributions are welcome! Please open issues or submit pull requests if you have any ideas, bug fixes, or enhancements.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
