package schema

import (
	"github.com/papaya147/go-cassandra-codegen/util"
)

const keyspaceNamePatternFromCreateKeyspace = `create keyspace ([a-z_0-9]+)`

func (schema *Schema) loadKeyspaces(ddl ...string) {
	for _, stmt := range ddl {
		if !util.CheckMatch(keyspaceNamePatternFromCreateKeyspace, stmt) {
			continue
		}

		keyspaceName := getMatch(keyspaceNamePatternFromCreateKeyspace, stmt)

		schema.Keyspaces = append(schema.Keyspaces, keyspaceName)
	}
}
