package query

import (
	"github.com/papaya147/go-cassandra-codegen/schema"
	"github.com/papaya147/go-cassandra-codegen/util"
)

var errorList = util.NewErrorList()

type argMap struct {
	Argument string
	ArgType  ArgType
}

type queryDetails struct {
	Keyspace     string
	Table        string
	Arguments    []argMap
	ReturnFields []string
	ReturnAmount QueryReturnAmount
}

type QueryList struct {
	schema *schema.Schema
	Query  map[string]queryDetails
}

func LoadQuery(schema *schema.Schema) (*QueryList, error) {
	queries := QueryList{
		schema: schema,
		Query:  map[string]queryDetails{},
	}

	files, err := util.GetFilesInDir(schema.Options.Cql.QueriesDir, "sql")
	if err != nil {
		return nil, err
	}

	dml, err := loadFiles(schema.Options.Cql.QueriesDir+"/", files...)
	if err != nil {
		return nil, err
	}

	queries.loadQueries(dml...)

	if !errorList.IsEmpty() {
		return nil, errorList.SerialiseError()
	}

	return &queries, nil
}
