package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/konjoot/kaboom/encoder"
)

func main() {
	var ruleString string
	flag.StringVar(ruleString, "rules", "", "semicolon separated string with encoding rules, for example 'one:string;two:int64'")
	flag.Parse()

	rules, err := encoder.ParseRules(rulesString)
	if err != nil {
		log.Println(err)
	}
	bts, err := encoder.Encode(os.Stdin, rules)
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
