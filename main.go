package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sort"

	"github.com/google/uuid"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	address string
	proxy   *httputil.ReverseProxy
}

type LoadBalancer struct {
	port    string
	servers []Server
}

// type ServerNode struct {
// 	Name string
// 	Host string
// }

type ConsistentHash struct {
	Keys       []int          // Sorted keys in the hash ring.
	Nodes      map[int]Server // Node mapping with hashed keys.
	TotalSlots int            // Total slots in the ring (hash space).
}

func NewConsistentHash(totalSlots int) *ConsistentHash {
	return &ConsistentHash{
		Nodes:      make(map[int]Server),
		TotalSlots: totalSlots,
	}
}

func hashFn(key string, totalSlots int) int {
	h := sha256.New()
	h.Write([]byte(key))
	hashValue := int(h.Sum(nil)[0]) // Using the first byte for simplicity
	return hashValue % totalSlots
}

func (ch *ConsistentHash) AddNode(nodeAddr Server) int {
	key := hashFn(nodeAddr.Address(), ch.TotalSlots)
	// Avoid collision by finding a new slot if the key is already taken.
	for {
		if _, exists := ch.Nodes[key]; !exists {
			break
		}
		key = (key + 1) % ch.TotalSlots
	}
	ch.Nodes[key] = nodeAddr
	ch.Keys = append(ch.Keys, key)
	sort.Ints(ch.Keys) // Keep the ring sorted
	return key
}

// RemoveNode removes a node from the hash ring.
func (ch *ConsistentHash) RemoveNode(nodeAddr Server) {
	key := hashFn(nodeAddr.Address(), ch.TotalSlots)
	index := sort.SearchInts(ch.Keys, key)
	if index < len(ch.Keys) && ch.Keys[index] == key {
		ch.Keys = append(ch.Keys[:index], ch.Keys[index+1:]...)
		delete(ch.Nodes, key)
	}
}

func (ch *ConsistentHash) Assign(item string) Server {
	key := hashFn(item, ch.TotalSlots)
	// Find the first node to the right of the hash
	index := sort.SearchInts(ch.Keys, key)
	if index == len(ch.Keys) { // Wrap around if necessary
		index = 0
	}
	return ch.Nodes[ch.Keys[index]]
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error:", err)
		os.Exit(1)
	}
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	handleError(err)

	return &simpleServer{
		address: addr,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {

	return &LoadBalancer{
		port:    port,
		servers: servers,
	}
}

func (s *simpleServer) Address() string {
	return s.address
}

func (s *simpleServer) IsAlive() bool {
	return true
}

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

func (lb *LoadBalancer) getNextAvailableServer(ch *ConsistentHash) Server {

	server := ch.Assign(uuid.New().String())

	return server // return the available server instance
}

func (lb *LoadBalancer) serverProxy(rw http.ResponseWriter, r *http.Request, ch *ConsistentHash) {
	targetServer := lb.getNextAvailableServer(ch)
	fmt.Printf("Forwarding request to address %s\n", targetServer.Address())
	targetServer.Serve(rw, r)
}

func main() {
	servers := []Server{
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("https://www.google.com"),
		newSimpleServer("https://www.instagram.com"),
	}

	lb := NewLoadBalancer("8000", servers)

	ch := NewConsistentHash(50)
	for _, node := range servers {
		key := ch.AddNode(node)
		fmt.Printf("Added node %s at key %d\n", node.Address(), key)
	}

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serverProxy(rw, req, ch)
	}

	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Starting server on port %s\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
