package main

import (
	"log"
	"net"

	"github.com/konjoot/kaboom/config"
	"github.com/konjoot/kaboom/mock"
	grpc "google.golang.org/grpc"
)

func main() {
	conf := config.New()

	lis, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		log.Println(err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	mock.RegisterMockServer(s, &mock.Endpoint{})

	s.Serve(lis)
}
