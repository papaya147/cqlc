package main

import (
	"fmt"
	"log"

	"github.com/papaya147/go-cassandra-codegen/options"
)

func main() {
	opts, err := options.NewOptions()
	if err != nil {
		fmt.Println(err)
	}

	log.Println(opts)
}
