package main

import (
	"log"

	"github.com/papaya147/go-cassandra-codegen/options"
)

func main() {
	opts, err := options.NewOptions("./codegen.yaml")
	if err != nil {
		log.Panic(err)
	}

	log.Println(*opts)
}
