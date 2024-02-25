package options

func NewMockOptions() *Options {
	opts := &Options{
		Version: 1,
		Cql: cqlConfig{
			Package:    "db",
			OutDir:     "db/cassandra-codegen",
			CqlPackage: "github.com/gocql/gocql",
			QueriesDir: "./db/queries",
			SchemaDir:  "./db/schema",
			Overrides:  nil,
		},
		Dependencies: []string{
			"time.Time",
			"string",
			"github.com/google/uuid.UUID",
		},
	}

	opts.replaceTypes()

	return opts
}
