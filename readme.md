# Load Balancer with Go Using Consistent Hashing

This project demonstrates a simple load balancer in Go that uses consistent hashing to distribute incoming HTTP requests across multiple backend servers. Consistent hashing ensures a relatively even distribution of load among servers and reduces disruption when servers are added or removed.

## Features

- **Consistent Hashing**: Uses a consistent hashing algorithm to map requests to servers, reducing cache misses when nodes are added or removed.
- **Reverse Proxy**: Uses `httputil.ReverseProxy` to forward incoming requests to backend servers.
- **Dynamic Node Assignment**: Backend servers can be dynamically added or removed from the consistent hash ring.

## Overview

The implementation uses `github.com/google/uuid` for generating unique keys for each incoming request. Based on this key, the consistent hash function determines the nearest server node to forward the request. This way, requests are evenly distributed, and hash collisions are managed.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/consistent-hashing-lb.git
   cd consistent-hashing-lb
   ```

2. Initialize and download dependencies:
   ```bash
   go mod tidy
   ```

3. Install required dependencies:
   ```bash
   go get github.com/google/uuid
   ```

## Structure

### Core Components

- **`Server` interface**: Defines basic methods (`Address`, `IsAlive`, `Serve`) for each backend server.
- **`simpleServer` struct**: Implements the `Server` interface and uses `httputil.ReverseProxy` to forward requests.
- **`LoadBalancer` struct**: Manages the list of servers and handles request distribution.
- **`ConsistentHash` struct**: Manages the hash ring and mapping of request keys to server nodes.

### Important Methods

1. **`AddNode`**: Adds a server to the hash ring, assigning it a hash key based on its address.
2. **`RemoveNode`**: Removes a server from the hash ring, deleting its key mapping.
3. **`Assign`**: Finds the appropriate server node for a given request key by looking up the closest server in the hash ring.
4. **`serverProxy`**: Routes incoming requests to the server selected by the consistent hashing mechanism.

## Usage

1. **Define Servers**: In `main()`, initialize backend servers using the `newSimpleServer` function.
2. **Set Up Load Balancer**: Initialize a `LoadBalancer` instance with the port number and list of servers.
3. **Set Up Consistent Hashing**: Create a `ConsistentHash` instance and add the server nodes to the hash ring.
4. **Start HTTP Server**: Set up an HTTP server that forwards requests to the load balancer’s `serverProxy` function, and listens for incoming requests.

## Example

To run the load balancer, simply execute:

```bash
go run main.go
```

This starts a server on `localhost:8000`, which will forward requests to either Facebook, Google, or Instagram, as per the consistent hashing algorithm.

### Sample Output

```
Added node https://www.facebook.com at key 4
Added node https://www.google.com at key 22
Added node https://www.instagram.com at key 30
Starting server on port 8000
Forwarding request to address https://www.google.com
Forwarding request to address https://www.facebook.com
```

### Test the Load Balancer

Once the server is running, make a few HTTP requests to `http://localhost:8000` using a browser or a tool like `curl`:

```bash
curl http://localhost:8000
```

You should see requests being forwarded to different servers based on the consistent hashing logic.

## Dependencies

- **Go Modules**: Ensure you have Go modules enabled.
- **External Packages**:
  - `github.com/google/uuid` for generating unique request IDs.
  - `net/http/httputil` for reverse proxying requests.
  - `crypto/sha256` for hashing.

This implementation provides a foundational understanding of consistent hashing-based load balancing, reverse proxying, and basic server health management.