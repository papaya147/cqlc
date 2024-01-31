package schema

import (
	"fmt"
	"testing"

	"github.com/papaya147/go-cassandra-codegen/util"
	"github.com/stretchr/testify/require"
)

func addRandomTableToSchema(t *testing.T) {
	keyspace := addRandomKeyspaceToSchema(t)
	table := util.RandomString(10)
	stmt := fmt.Sprintf(`create table %s.%s(
		t1 int,
		t2 int,
		primary key((t1), t2)
	)`, keyspace, table)
	testSchema.loadTables(stmt)
	require.Empty(t, errorList)
	require.Contains(t, testSchema.Tables, Table{
		Keyspace: keyspace,
		Name:     table,
		Fields: map[string]string{
			"t1": "int32",
			"t2": "int32",
		},
		PartitionKeys:  []string{"t1"},
		ClusteringKeys: []string{"t2"},
	})
}

func TestLoadTables(t *testing.T) {
	wipeTestSchema()
	addRandomTableToSchema(t)
}
