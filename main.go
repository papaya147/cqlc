package main

import (
	"fmt"

	"github.com/papaya147/go-cassandra-codegen/options"
)

func main() {
	opts, err := options.NewOptions()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(opts)
}
