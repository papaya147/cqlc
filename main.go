package main

import (
	"fmt"

	"github.com/papaya147/go-cassandra-codegen/options"
	"github.com/papaya147/go-cassandra-codegen/schema"
)

func main() {
	opts, err := options.NewOptions()
	if err != nil {
		fmt.Println(err)
	}

	s, err := schema.LoadSchema(opts)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(s)
}
