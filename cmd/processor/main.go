package main

import (
	"flag"
	"log"
	"os"

	"github.com/konjoot/kaboom/processor"
)

func main() {

	var addr, method string
	flag.StringVar(&addr, "addr", "", "address of the gRPC endpoint")
	flag.StringVar(&method, "method", "", "method name of the gRPC endpoint")
	flag.Parse()

	err := processor.Process(os.Stdin, addr, method, os.Stdout)
	if err != nil {
		log.Println(err)
	}
}
