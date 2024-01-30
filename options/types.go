package options

import (
	"fmt"
	"strings"
)

var DefaultTypeMapping = map[string]string{
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

func (opts *Options) replaceTypes() error {
	typesToReplace := opts.Cql.Overrides
	for _, v := range typesToReplace {
		if _, ok := DefaultTypeMapping[v.DbType]; !ok {
			return fmt.Errorf("type %s not supported or does not exist", v.DbType)
		}
		goTypeSplit := strings.Split(v.GoType, "/")
		goType := goTypeSplit[len(goTypeSplit)-1]
		DefaultTypeMapping[v.DbType] = goType

		if strings.Contains(goType, ".") {
			goPackageSplit := strings.Split(v.GoType, ".")
			goPackageSplit = goPackageSplit[:len(goPackageSplit)-1]
			opts.Dependencies = append(opts.Dependencies, strings.Join(goPackageSplit, "."))
		}
	}

	opts.Dependencies = append(opts.Dependencies, opts.Cql.CqlPackage)
	opts.TypeMappings = DefaultTypeMapping
	return nil
}
