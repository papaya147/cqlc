package ddl

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/papaya147/cqlc/util"
)

var rawStatements []string

func Load(ctx context.Context, dir string) error {
	stmts, err := getDDL(ctx, dir)
	if err != nil {
		return err
	}

	rawStatements = stmts

	// loading keyspaces and tables
	loadKeyspaces(ctx)
	if err := loadTables(ctx); err != nil {
		return err
	}

	return nil
}

func getDDL(ctx context.Context, dir string) ([]string, error) {
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
