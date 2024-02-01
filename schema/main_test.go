package schema

import (
	"os"
	"testing"

	"github.com/papaya147/go-cassandra-codegen/options"
)

var testSchema = Schema{}

func wipeTestSchema() {
	testSchema = Schema{
		Options: options.NewMockOptions(),
		Tables:  map[string]table{},
	}
}

func TestMain(m *testing.M) {
	testSchema.Options = options.NewMockOptions()
	testSchema.Tables = map[string]table{}

	os.Exit(m.Run())
}
