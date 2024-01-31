package schema

import (
	"fmt"
	"testing"

	"github.com/papaya147/go-cassandra-codegen/util"
	"github.com/stretchr/testify/require"
)

func addRandomKeyspaceToSchema(t *testing.T) string {
	keyspace := util.RandomString(10)
	stmt := fmt.Sprintf("create keyspace %s", keyspace)
	testSchema.loadKeyspaces(stmt)
	require.NotEmpty(t, testSchema.Keyspaces)
	require.Empty(t, errorList)
	require.Contains(t, testSchema.Keyspaces, keyspace)

	return keyspace
}

func TestLoadKeyspaces(t *testing.T) {
	wipeTestSchema()
	addRandomKeyspaceToSchema(t)
}
