package main

import (
	"context"
	"log"
	protos "main/protos"

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

	client := protos.NewServiceClient(conn)

	data, err := client.GetData(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic(err)
	}

	log.Printf("Name: %s, Value: %s\n", data.Name, data.Value)
}
