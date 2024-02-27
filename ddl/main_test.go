package ddl

import (
	"os"
	"testing"

	"github.com/papaya147/cqlc/options"
)

func TestMain(m *testing.M) {
	userOptions = options.NewMockOptions()
	os.Exit(m.Run())
}
