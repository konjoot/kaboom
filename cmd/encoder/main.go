package main

import (
	"fmt"
	"log"
	"os"

	"github.com/konjoot/kaboom/encoder"
)

func main() {
	bts, err := encoder.Encode(os.Stdin)
	if err != nil {
		log.Println(err)
	}
	n, err := fmt.Fprintln(os.Stdout, bts)
	if err != nil {
		log.Println(err)
	}
	if n < len(bts) {
		log.Println("the number of bytes written to STDOUT is less than expected")
	}
}
