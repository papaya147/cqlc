package schema

import (
	"github.com/papaya147/go-cassandra-codegen/util"
)

const keyspaceNamePatternFromCreateKeyspace = `create keyspace ([a-z_]+)`

func (schema *Schema) loadKeyspaces(ddl ...string) error {
	for _, stmt := range ddl {
		if !util.CheckMatch(keyspaceNamePatternFromCreateKeyspace, stmt) {
			continue
		}

		keyspaceName := schema.getKeyspace(stmt)

		schema.Keyspaces = append(schema.Keyspaces, keyspaceName)
	}

	return nil
}

func (schema *Schema) getKeyspace(stmt string) string {
	match, err := util.GetFirstMatch(keyspaceNamePatternFromCreateKeyspace, stmt)
	if err != nil {
		errorList.Add(err)
	}

	return match
}
