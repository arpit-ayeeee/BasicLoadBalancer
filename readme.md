# Basic Load Balancer

This Go package implements a simple HTTP load balancer using round-robin scheduling. The load balancer forwards incoming HTTP requests to a list of backend servers in a circular manner. If a server is not responding, the load balancer skips to the next available server.

## Features

- **Round-robin load balancing**: Distributes requests evenly across backend servers.
- **Health Check**: Skips non-responsive servers and continues routing to live servers.
- **Reverse Proxy**: Proxies incoming requests to target servers.

## Requirements

- Go 1.16 or higher

## Getting Started

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-repo-url.git
   ```

2. **Navigate to the project directory**:
   ```bash
   cd your-project-directory
   ```

### Usage

The package creates a load balancer that listens on a specified port (default `8000` in this code) and forwards requests to a list of specified backend servers.

### Code Walkthrough

- **main()**: Initializes the server instances and starts the load balancer on port `8000`.
- **LoadBalancer**:
  - `getNextAvailableServer()`: Selects the next available server using a round-robin algorithm.
  - `serverProxy()`: Routes the HTTP request to the selected server.
- **simpleServer**:
  - `Serve()`: Handles request forwarding using the `httputil.ReverseProxy`.
  - `IsAlive()`: Currently returns `true` for all servers (can be extended for actual health checks).

### Running the Code

To run the load balancer:
```bash
go run main.go
```

### Example Configuration

The `main.go` initializes three backend servers:

```go
servers := []Server{
    newSimpleServer("https://www.facebook.com"),
    newSimpleServer("https://www.google.com"),
    newSimpleServer("https://www.instagram.com"),
}
```

The server listens on port `8000` and routes requests to the backend servers in a round-robin sequence.

### Load Balancer Struct

- **Fields**:
  - `port`: The port on which the load balancer listens.
  - `roundRobinCount`: Keeps track of the last served server for round-robin selection.
  - `servers`: A slice of `Server` interfaces representing backend servers.

## API Endpoints

- **/**: All incoming requests are routed through the load balancer, which selects the next available server and proxies the request.

## Example Output

Upon running, the server displays logs of incoming requests and the server addresses to which requests are forwarded:

```plaintext
Starting server on port 8000
Forwarding request to address https://www.facebook.com
Forwarding request to address https://www.google.com
Forwarding request to address https://www.instagram.com
```

## Extending the Package

1. **Health Checks**: Modify `IsAlive()` to perform an actual health check to ensure the backend server is available.
2. **Dynamic Server Addition**: Add functionality to add or remove servers dynamically.
