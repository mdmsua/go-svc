package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var infoLog = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
var warnLog = log.New(os.Stdout, "[WARN] ", log.LstdFlags)
var erroLog = log.New(os.Stderr, "[ERRO] ", log.LstdFlags)

type Server struct {
	hostname string
	cert     string
	key      string
	port     int
	grpc     int
}

type service struct {
	UnimplementedServiceServer
}

func (service *service) GetData(ctx context.Context, req *emptypb.Empty) (*Data, error) {
	return &Data{
		Name:  getHostname(),
		Value: "Hello World",
	}, nil
}

func NewServer(cert string, key string, port int, grpc int) *Server {
	return &Server{
		hostname: getHostname(),
		cert:     cert,
		key:      key,
		port:     port,
		grpc:     grpc,
	}
}

func (s Server) Run() {
	infoLog.Printf("Starting server on port %d, gRPC port %d, TLS cert path: %s, TLS key path: %s\n", s.port, s.grpc, s.cert, s.key)

	var data Response

	go func() {
		client := &http.Client{
			Transport: &http.Transport{},
		}

		resp, err := client.Get("https://ifconfig.me/all.json")
		if err != nil {
			erroLog.Printf("Failed to get response: %v\n", err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			erroLog.Printf("Failed to read response body: %v\n", err)
			return
		}

		if json.Unmarshal(body, &data) != nil {
			erroLog.Printf("Failed to unmarshal response body: %v\n", err)
			return
		}
	}()

	http.HandleFunc("GET /data", func(w http.ResponseWriter, r *http.Request) {
		infoLog.Printf("Handling request from %v\n", r.RemoteAddr)
		w.Write([]byte(fmt.Sprintf("Hello World from %v to %v via %v, egress %v via %v", s.hostname, r.RemoteAddr, r.Proto, data.egress(), data.Fowarded)))
	})

	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	httpAddr := fmt.Sprintf("0.0.0.0:%d", s.port)
	grpcAddr := fmt.Sprintf("0.0.0.0:%d", s.grpc)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		erroLog.Fatalf("(gRPC) failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	RegisterServiceServer(grpcServer, &service{})
	go func() {
		grpcServer.Serve(lis)
	}()

	if s.cert != "" && s.key != "" {
		erroLog.Fatal(http.ListenAndServeTLS(httpAddr, s.cert, s.key, nil))
	} else {
		erroLog.Fatal(http.ListenAndServe(httpAddr, nil))
	}
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		warnLog.Printf("failed to get hostname: %v", err)
		hostname = "N/A"
	}

	return hostname
}
