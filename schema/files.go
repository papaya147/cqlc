package schema

import (
	"fmt"
	"strings"

	"github.com/papaya147/go-cassandra-codegen/util"
)

func loadFiles(prefix string, files ...string) ([]string, error) {
	result := []string{}
	for _, file := range files {
		fileContents, err := util.GetFileContents(fmt.Sprintf("%s%s", prefix, file))
		if err != nil {
			return nil, err
		}

		ddlSplit := strings.Split(string(fileContents), ";")
		for _, ddl := range ddlSplit {
			q := strings.ToLower(strings.TrimSpace(ddl))
			if ddl != "" {
				result = append(result, q)
			}
		}
	}

	return result, nil
}
