package ddl

import (
	"context"
	"fmt"
	"strings"

	"github.com/papaya147/cqlc/util"
)

var keyspaces []string
var tables []string

const captureKeyspaceFromCreateKeyspace = `(?i)CREATE\s+KEYSPACE\s+([A-Za-z_][A-Za-z0-9_]*)`

func loadKeyspaceNames(context.Context) {
	for _, stmt := range rawStatements {
		keyspace, err := util.GetFirstMatch(captureKeyspaceFromCreateKeyspace, stmt)
		if err == nil {
			keyspaces = append(keyspaces, strings.ToLower(keyspace))
		}
	}
}

var tableKeyspaceMap map[string]string

const captureKeyspaceFromCreateTable = `(?i)CREATE\s+TABLE\s+([A-Za-z_][A-Za-z_0-9]*)\.`
const captureTableFromCreateTable = `(?i)CREATE\s+TABLE[^.]+\.([A-Za-z_][A-Za-z_0-9]*)`

func loadTableNames(context.Context) error {
	list := util.NewErrorList()
	tableKeyspaceMap = make(map[string]string)

	for _, stmt := range rawStatements {
		// getting keyspace in table definition
		keyspace, err := util.GetFirstMatch(captureKeyspaceFromCreateTable, stmt)
		if err != nil {
			continue
		}
		keyspace = strings.ToLower(keyspace)

		if !util.Contains(keyspaces, keyspace) {
			list.Add(fmt.Errorf("keyspace %s was not created", keyspace))
		}

		table, err := util.GetFirstMatch(captureTableFromCreateTable, stmt)
		if err == nil {
			table = strings.ToLower(table)
			tables = append(tables, table)
			tableKeyspaceMap[table] = keyspace
		}
	}

	if list.IsEmpty() {
		return nil
	}

	return list
}

type Fields map[string]string

var tableFields map[string]Fields

type Keys struct {
	PartitionKeys  []string
	ClusteringKeys []string
}

var tableKeys map[string]Keys

const captureTableFieldsFromCreateTable = `(?i)([A-Za-z_0-9]+\s+[A-Za-z_0-9]+),`
const capturePartitionKeysFromCreateTable = `(?i)PRIMARY\s+KEY\s*\(\(([A-Za-z_0-9\s,]+)\)`
const captureClusteringKeysFromCreateTable = `(?i)PRIMARY\s+KEY\s*\(\([^)]+\),\s*([A-Za-z_0-9]+)`

func loadTableFields(ctx context.Context, supportedTypes map[string]string) error {
	list := util.NewErrorList()
	tableFields = make(map[string]Fields)
	tableKeys = make(map[string]Keys)

	for _, stmt := range rawStatements {
		// getting table in table definition
		table, err := util.GetFirstMatch(captureTableFromCreateTable, stmt)
		if err != nil {
			continue
		}
		table = strings.ToLower(table)

		// assigning table fields
		fields, err := util.GetAllMatches(captureTableFieldsFromCreateTable, stmt)
		if err != nil || len(fields) == 0 {
			list.Add(fmt.Errorf("table %s has no fields", table))
		}

		if tableFields[table] != nil {
			list.Add(fmt.Errorf("table %s has duplicate definitions", table))
		}

		tableFields[table] = make(Fields)

		for _, f := range fields {
			fieldAndType := strings.Split(f, " ")
			if len(fieldAndType) < 2 {
				list.Add(fmt.Errorf("field %s has no type", f))
			}

			field := strings.ToLower(fieldAndType[0])
			fieldType := strings.ToLower(fieldAndType[1])

			if supportedTypes[fieldType] == "" {
				list.Add(fmt.Errorf("field %s has unsupported type %s", field, fieldType))
			}

			if tableFields[table][field] != "" {
				list.Add(fmt.Errorf("field %s has duplicate definition in table %s", field, table))
			}

			tableFields[table][field] = fieldType
		}

		// assigning table keys
		k := Keys{}

		partitionKeyBlock, err := util.GetFirstMatch(capturePartitionKeysFromCreateTable, stmt)
		if err != nil {
			list.Add(fmt.Errorf("table %s has no partition key", table))
		}

		partitionKeys := strings.Split(partitionKeyBlock, ",")
		if len(partitionKeys) == 0 {
			list.Add(fmt.Errorf("table %s has no partition key", table))
		}

		for _, p := range partitionKeys {
			p = strings.ToLower(strings.TrimSpace(p))
			if tableFields[table][p] == "" {
				list.Add(fmt.Errorf("partition key %s in table %s is not defined", p, table))
			}

			k.PartitionKeys = append(k.PartitionKeys, p)
		}

		clusteringKeyBlock, err := util.GetFirstMatch(captureClusteringKeysFromCreateTable, stmt)
		if err != nil {
			list.Add(fmt.Errorf("table %s has no clustering key", table))
		}

		clusteringKeys := strings.Split(clusteringKeyBlock, ", ")
		if len(clusteringKeys) == 0 {
			list.Add(fmt.Errorf("table %s has no clustering key", table))
		}

		for _, c := range clusteringKeys {
			c = strings.ToLower(strings.TrimSpace(c))
			if tableFields[table][c] == "" {
				list.Add(fmt.Errorf("clustering key %s in table %s is not defined", c, table))
			}

			k.ClusteringKeys = append(k.ClusteringKeys, c)
		}

		tableKeys[table] = k
	}

	if list.IsEmpty() {
		return nil
	}

	return list
}
