package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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
	port            string
	roundRobinCount int
	servers         []Server
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
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
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

func (lb *LoadBalancer) getNextAvailableServer() Server {

	server := lb.servers[lb.roundRobinCount%len(lb.servers)]

	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++

	return server // return the available server instance
}

func (lb *LoadBalancer) serverProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
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

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serverProxy(rw, req)
	}

	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Starting server on port %s\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
