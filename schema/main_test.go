package schema

import (
	"os"
	"testing"

	"github.com/papaya147/go-cassandra-codegen/options"
)

var testSchema = Schema{}

func wipeTestSchema() {
	testSchema = Schema{
		options: options.NewMockOptions(),
	}
}

func TestMain(m *testing.M) {
	testSchema.options = options.NewMockOptions()

	os.Exit(m.Run())
}
