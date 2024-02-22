package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	svc "main/services"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	cert string
	key  string
	port int
	grpc int
}

type service struct {
	svc.UnimplementedServiceServer
}

func (service *service) GetData(ctx context.Context, req *emptypb.Empty) (*svc.Data, error) {
	return &svc.Data{
		Timestamp: timestamppb.Now(),
		Name:      "Server",
		Value:     "Hello World",
	}, nil
}

func NewServer(cert string, key string, port int, grpc int) *Server {
	return &Server{
		cert: cert,
		key:  key,
		port: port,
		grpc: grpc,
	}
}

func (s Server) Run() {
	infLog := log.New(os.Stdout, "[INF] ", log.LstdFlags)
	errLog := log.New(os.Stderr, "[ERR] ", log.LstdFlags)

	infLog.Printf("Starting server on port %d, gRPC port %d, TLS cert path: %s, TLS key path: %s\n", s.port, s.grpc, s.cert, s.key)

	hostname, err := os.Hostname()
	if err != nil {
		errLog.Printf("Failed to get hostname: %v\n", err)
	}

	var data Response

	go func() {
		client := &http.Client{
			Transport: &http.Transport{},
		}

		resp, err := client.Get("https://ifconfig.me/all.json")
		if err != nil {
			errLog.Printf("Failed to get response: %v\n", err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			errLog.Printf("Failed to read response body: %v\n", err)
		}

		if json.Unmarshal(body, &data) != nil {
			errLog.Printf("Failed to unmarshal response body: %v\n", err)
		}
	}()

	http.HandleFunc("GET /data", func(w http.ResponseWriter, r *http.Request) {
		infLog.Printf("Handling request from %v\n", r.RemoteAddr)
		w.Write([]byte(fmt.Sprintf("Hello World from %v to %v via %v, egress %v via %v", hostname, r.RemoteAddr, r.Proto, data.egress(), data.Fowarded)))
	})

	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	httpAddr := fmt.Sprintf("0.0.0.0:%d", s.port)
	grpcAddr := fmt.Sprintf("0.0.0.0:%d", s.grpc)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		errLog.Fatalf("(gRPC) failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	svc.RegisterServiceServer(grpcServer, &service{})
	go func() {
		grpcServer.Serve(lis)
	}()

	if s.cert != "" && s.key != "" {
		errLog.Fatal(http.ListenAndServeTLS(httpAddr, s.cert, s.key, nil))
	} else {
		errLog.Fatal(http.ListenAndServe(httpAddr, nil))
	}
}
