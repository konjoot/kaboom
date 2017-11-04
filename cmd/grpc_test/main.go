package main

import (
	"context"
	"log"
	"os"

	"github.com/konjoot/kaboom/config"
	"github.com/konjoot/kaboom/mock"
	"google.golang.org/grpc"
)

func main() {

	conf := config.New()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	conn, err := grpc.Dial(conf.Listen, opts...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	cli := mock.NewMockClient(conn)
	_, err = cli.Base(context.Background(), &mock.BaseMsg{Int32: -150, Int64: -150})
	log.Println(err)
}
