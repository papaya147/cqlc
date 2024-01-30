package schema

import (
	"github.com/papaya147/go-cassandra-codegen/options"
	"github.com/papaya147/go-cassandra-codegen/util"
)

type Table struct {
	Keyspace       string
	Name           string
	Fields         map[string]string
	PartitionKeys  []string
	ClusteringKeys []string
}

type Schema struct {
	Tables    []Table
	Keyspaces []string
	options   *options.Options
}

var errorList = util.NewErrorList()

func LoadSchema(opts *options.Options) (*Schema, error) {
	schema := Schema{
		options: opts,
	}

	files, err := util.GetFilesInDir(schema.options.Cql.SchemaDir, "sql")
	if err != nil {
		return nil, err
	}

	ddl, err := loadFiles(schema.options.Cql.SchemaDir+"/", files...)
	if err != nil {
		return nil, err
	}

	if err := schema.loadKeyspaces(ddl...); err != nil {
		return nil, err
	}

	schema.loadTables(ddl...)

	if !errorList.IsEmpty() {
		return nil, errorList.SerialiseError()
	}

	return &schema, nil
}
