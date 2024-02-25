package main

import (
	"fmt"

	"github.com/papaya147/cqlc/options"
)

func main() {
	opts, err := options.LoadOptions()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(opts)
}
