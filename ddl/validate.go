package ddl

import (
	"context"
	"fmt"

	"github.com/papaya147/cqlc/util"
)

var keyspaces []string
var tables []string

const captureKeyspaceFromCreateKeyspace = `(?i)CREATE KEYSPACE ([A-Za-z_][A-Za-z0-9_]*)`

func loadKeyspaces(context.Context) {
	for _, stmt := range rawStatements {
		keyspace, err := util.GetFirstMatch(captureKeyspaceFromCreateKeyspace, stmt)
		if err == nil {
			keyspaces = append(keyspaces, keyspace)
		}
	}
}

const captureKeyspaceFromCreateTable = `(?i)CREATE TABLE ([A-Za-z_][A-Za-z_0-9]*)\.`
const captureTableFromCreateTable = `(?i)CREATE TABLE [^.]+\.([A-Za-z_][A-Za-z_0-9]*)`

func loadTables(context.Context) error {
	list := util.NewErrorList()

	for _, stmt := range rawStatements {
		// getting keyspace in table definition
		keyspace, err := util.GetFirstMatch(captureKeyspaceFromCreateTable, stmt)
		if err != nil {
			continue
		}

		if !util.Contains(keyspaces, keyspace) {
			list.Add(fmt.Errorf("keyspace %s was not created", keyspace))
		}

		table, err := util.GetFirstMatch(captureTableFromCreateTable, stmt)
		if err == nil {
			tables = append(tables, table)
		}
	}

	if list.IsEmpty() {
		return nil
	}

	return list
}
