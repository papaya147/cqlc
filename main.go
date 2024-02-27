package main

import (
	"context"
	"fmt"
	"log"

	"github.com/papaya147/cqlc/ddl"
	"github.com/papaya147/cqlc/dml"
	"github.com/papaya147/cqlc/options"
	"github.com/papaya147/parallelize"
)

func main() {
	opts, err := options.LoadOptions()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := loadFiles(opts); err != nil {
		fmt.Println(err)
		return
	}

	ddlConfig, err := ddl.PrepareConfig(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	dmlConfig, err := dml.PrepareConfig(context.Background(), ddlConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, a := range dmlConfig {
		log.Println(a.Name)
		for k, v := range a.InputFields {
			log.Println(k)
			log.Println(*v)
		}
		log.Println()
	}

	fmt.Println("no errors")
}

func loadFiles(opts *options.Config) error {
	group := parallelize.NewSyncGroup()

	parallelize.AddMethodWithArgs(group, ddl.Load, parallelize.MethodWithArgsParams[*options.Config]{
		Context: context.Background(),
		Input:   opts,
	})

	parallelize.AddMethodWithArgs(group, dml.Load, parallelize.MethodWithArgsParams[*options.Config]{
		Context: context.Background(),
		Input:   opts,
	})

	return group.Run()
}
