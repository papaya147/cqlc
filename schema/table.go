package schema

import (
	"fmt"
	"strings"

	"github.com/papaya147/go-cassandra-codegen/util"
)

func (schema *Schema) loadTables(ddl ...string) {
	for _, stmt := range ddl {
		if !util.CheckMatch(tableNamePatternFromCreateTable, stmt) {
			continue
		}

		schema.getTable(stmt)
	}
}

func (schema *Schema) getTable(stmt string) {
	tableName := getMatch(tableNamePatternFromCreateTable, stmt)

	keyspace := getMatch(keyspaceNamePatternFromCreateTable, stmt)

	// check if keyspace exists in create statements
	if !util.Contains[string](schema.Keyspaces, keyspace) {
		errorList.Add(fmt.Errorf("keyspace %s not found in create statements", keyspace))
	}

	fields := getFields(stmt)

	partitionKeys := getPartitionKeys(stmt)

	// check if partition keys in fields
	for _, key := range partitionKeys {
		if _, ok := fields[key]; !ok {
			errorList.Add(fmt.Errorf("partition key %s not found in fields", key))
		}
	}

	clusteringKeys := getClusteringKeys(stmt)

	// check if clustering keys in fields
	for _, key := range clusteringKeys {
		if _, ok := fields[key]; !ok {
			errorList.Add(fmt.Errorf("clustering key %s not found in fields", key))
		}
	}

	updatedTypeFields := schema.injectOptionTypes(fields)

	schema.Tables[tableName] = table{
		Keyspace:       keyspace,
		Fields:         updatedTypeFields,
		PartitionKeys:  partitionKeys,
		ClusteringKeys: clusteringKeys,
	}
}

const tableNamePatternFromCreateTable = `create table [a-z_0-9]+\.([a-z_0-9]+)`

const keyspaceNamePatternFromCreateTable = `create table ([a-z_0-9]+)`

const fieldsPatternFromCreateTable = `create table [a-z0-9.]+\(([a-z_\s,0-9]+)primary key`

func getFields(stmt string) map[string]string {
	match := getMatch(fieldsPatternFromCreateTable, stmt)

	fieldsUntrimmed := strings.Split(match, ",")
	fields := make(map[string]string, len(fieldsUntrimmed))
	for _, field := range fieldsUntrimmed {
		fieldParts := strings.Split(strings.TrimSpace(field), " ")
		if len(fieldParts) < 2 {
			continue
		}

		fields[strings.TrimSpace(fieldParts[0])] = strings.TrimSpace(fieldParts[1])
	}

	return fields
}

const partitionKeysPatternFromCreateTable = `primary key\s*\(\s*\(([a-z_,\s0-9]+)`

func getPartitionKeys(stmt string) []string {
	match := getMatch(partitionKeysPatternFromCreateTable, stmt)

	keysUntrimmed := strings.Split(match, ",")
	keys := make([]string, len(keysUntrimmed))
	for i, key := range keysUntrimmed {
		keys[i] = strings.TrimSpace(key)
	}

	return keys
}

const clusteringKeyPatternFromCreateTable = `primary key\s*\(\s*\([a-z_,\s0-9]+\),\s*([a-z_,\s0-9]+)`

func getClusteringKeys(stmt string) []string {
	match := getMatch(clusteringKeyPatternFromCreateTable, stmt)

	keysUntrimmed := strings.Split(strings.TrimSpace(match), ",")
	keys := make([]string, len(keysUntrimmed))
	for i, key := range keysUntrimmed {
		keys[i] = strings.TrimSpace(key)
	}

	return keys
}

func (schema *Schema) injectOptionTypes(fields map[string]string) map[string]string {
	goFields := map[string]string{}
	for name, typ := range fields {
		goType, ok := schema.Options.TypeMappings[typ]
		if !ok {
			errorList.Add(fmt.Errorf("type %s is not a valid type for field %s", typ, name))
		}

		goFields[name] = goType
	}

	return goFields
}
