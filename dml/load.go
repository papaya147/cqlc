package dml

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/papaya147/cqlc/util"
)

var rawStatements []string

func Load(ctx context.Context, dir string) error {
	stmts, err := getDML(ctx, dir)
	if err != nil {
		return err
	}

	rawStatements = stmts

	return nil
}

func getDML(ctx context.Context, dir string) ([]string, error) {
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
		return nil, errors.New("no content found in any query files")
	}

	return strings.Split(fileContents, ";"), nil
}
