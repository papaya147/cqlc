package main

import (
	"context"
	"fmt"

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

	fmt.Println("no errors")
}

func loadFiles(opts *options.Options) error {
	group := parallelize.NewSyncGroup()

	parallelize.AddMethodWithArgs(group, ddl.Load, parallelize.MethodWithArgsParams[string]{
		Context: context.Background(),
		Input:   opts.Cql.SchemaDir,
	})

	parallelize.AddMethodWithArgs(group, dml.Load, parallelize.MethodWithArgsParams[string]{
		Context: context.Background(),
		Input:   opts.Cql.QueriesDir,
	})

	return group.Run()
}
