package client

import (
	"context"
	"log"
	svc "main/services"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c Client) Run() {
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := svc.NewServiceClient(conn)

	data, err := client.GetData(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic(err)
	}

	log.Printf("Timestamp: %v, Name: %s, Value: %s\n", data.Timestamp, data.Name, data.Value)
}
