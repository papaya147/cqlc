package options

import (
	"fmt"
	"log"
)

var CassandraTypeMapping = map[string]string{
	"ascii":     "string",
	"bigint":    "int64",
	"blob":      "[]byte",
	"boolean":   "bool",
	"counter":   "int64",
	"date":      "time.Time",
	"double":    "float64",
	"float":     "float32",
	"int":       "int32",
	"inet":      "string",
	"text":      "string",
	"time":      "time.Time",
	"timestamp": "time.Time",
	"tinyint":   "int8",
	"smallint":  "int16",
	"varchar":   "string",
	"varint":    "int",
	"uuid":      "uuid.UUID",
}

func replaceTypes(opts Options) error {
	typesToReplace := opts.Cql.Overrides
	for _, v := range typesToReplace {
		if _, ok := CassandraTypeMapping[v.DbType]; !ok {
			return fmt.Errorf("type %s not supported or does not exist", v.DbType)
		}
		CassandraTypeMapping[v.DbType] = v.GoType
	}
	log.Println(CassandraTypeMapping)
	return nil
}
