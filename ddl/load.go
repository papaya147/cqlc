package ddl

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/papaya147/cqlc/util"
)

var RawStatements []string

func Load(ctx context.Context, dir string) error {
	stmts, err := getDDL(dir)
	if err != nil {
		return err
	}

	RawStatements = stmts

	return nil
}

func getDDL(dir string) ([]string, error) {
	fileNames, err := util.GetFilesInDir(dir, "sql")
	if err != nil {
		return nil, err
	}

	fileContents := ""
	for _, fileName := range fileNames {
		filePath := fmt.Sprintf("%s/%s", dir, fileName)
		content, err := util.GetFileContents(filePath)
		if err != nil {
			return nil, err
		}
		fileContents += string(content)
	}

	if fileContents == "" {
		return nil, errors.New("no content found in any schema files")
	}

	return strings.Split(fileContents, ";"), nil
}
