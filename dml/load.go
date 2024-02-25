package dml

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/papaya147/cqlc/options"
	"github.com/papaya147/cqlc/util"
)

var userOptions *options.Config

var rawStatements []string

type Config struct {
}

func Load(ctx context.Context, opts *options.Config) error {
	stmts, err := getDML(ctx, opts.Cql.QueriesDir)
	if err != nil {
		return err
	}

	userOptions = opts
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
