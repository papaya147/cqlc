package schema

import (
	"github.com/papaya147/go-cassandra-codegen/options"
	"github.com/papaya147/go-cassandra-codegen/util"
)

type table struct {
	Keyspace       string
	Fields         map[string]string
	PartitionKeys  []string
	ClusteringKeys []string
}

type Schema struct {
	Tables    map[string]table
	Keyspaces []string
	Options   *options.Options
}

var errorList = util.NewErrorList()

func LoadSchema(opts *options.Options) (*Schema, error) {
	schema := Schema{
		Options: opts,
		Tables:  map[string]table{},
	}

	files, err := util.GetFilesInDir(schema.Options.Cql.SchemaDir, "sql")
	if err != nil {
		return nil, err
	}

	ddl, err := loadFiles(schema.Options.Cql.SchemaDir+"/", files...)
	if err != nil {
		return nil, err
	}

	schema.loadKeyspaces(ddl...)
	schema.loadTables(ddl...)

	if !errorList.IsEmpty() {
		return nil, errorList.SerialiseError()
	}

	return &schema, nil
}
