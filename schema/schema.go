package schema

import "github.com/papaya147/go-cassandra-codegen/options"

type Field struct {
	Name string
	Type string
}

type Schema struct {
	Keyspace string
	Table    string
	Fields   []Field
}

func NewSchema(opts *options.Options) {

}
