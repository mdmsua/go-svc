package main

import (
	"client"
	"flag"
	"log"
	"os"
	"server"
)

func main() {
	serverFlagSet := flag.NewFlagSet("server", flag.ExitOnError)
	clientFlagSet := flag.NewFlagSet("client", flag.ExitOnError)

	cert := serverFlagSet.String("cert", "", "Path to server TLS certificate")
	key := serverFlagSet.String("key", "", "Path to server TLS key")
	port := serverFlagSet.Int("port", 8080, "Port to run server on")
	grpc := serverFlagSet.Int("grpc", 8081, "Port to run gRPC server on")
	addr := clientFlagSet.String("addr", "localhost:8081", "Address of the gRPC server")

	switch os.Args[1] {
	case "server":
		serverFlagSet.Parse(os.Args[2:])
		server := server.NewServer(*cert, *key, *port, *grpc)
		server.Run()
	case "client":
		clientFlagSet.Parse(os.Args[2:])
		client := client.NewClient(*addr)
		client.Run()
	default:
		log.Fatal("Invalid command")
	}
}
