package ddl

import (
	"context"
	"fmt"
	"testing"

	"github.com/papaya147/cqlc/util"
	"github.com/stretchr/testify/require"
)

func createRandomKeyspaceStatement(t *testing.T) string {
	keyspace := util.RandomString(10)
	rawStatements = []string{
		fmt.Sprintf(`CREATE KEYSPACE %s`, keyspace),
	}
	return keyspace
}

func createRandomTableStatement(t *testing.T) (string, string, []string) {
	keyspace := createRandomKeyspaceStatement(t)
	table := util.RandomString(10)
	field1 := util.RandomString(10)
	field2 := util.RandomString(10)
	rawStatements = append(rawStatements, fmt.Sprintf(`
		CREATE TABLE %s.%s (
			%s varchar,
			%s bigint,
			PRIMARY KEY ((%s), %s)
		)
	`, keyspace, table, field1, field2, field1, field2))
	return keyspace, table, []string{field1, field2}
}

func TestLoadKeyspaceNames(t *testing.T) {
	keyspace1 := createRandomKeyspaceStatement(t)

	loadKeyspaceNames(context.Background())

	require.Contains(t, keyspaces, keyspace1)
}

func TestLoadTables(t *testing.T) {
	_, table, _ := createRandomTableStatement(t)

	loadKeyspaceNames(context.Background())
	err := loadTableNames(context.Background())
	require.NoError(t, err)

	require.Contains(t, tables, table)
}

func TestLoadTableFields(t *testing.T) {
	_, table, fieldNames := createRandomTableStatement(t)

	loadKeyspaceNames(context.Background())
	err := loadTableNames(context.Background())
	require.NoError(t, err)
	err = loadTableFields(context.Background(), userOptions.TypeMappings)
	require.NoError(t, err)

	require.Equal(t, tableFields[table][fieldNames[0]], "varchar")
	require.Equal(t, tableFields[table][fieldNames[1]], "bigint")
}
