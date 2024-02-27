package dml

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/papaya147/cqlc/ddl"
	"github.com/papaya147/cqlc/options"
	"github.com/papaya147/cqlc/util"
)

var userOptions *options.Config

var rawStatements []string

type QueryReturnAmount int

const (
	RETURN_NONE QueryReturnAmount = iota
	RETURN_ONE
	RETURN_MANY
)

type InputAmount int

const (
	INPUT_ONE InputAmount = iota + 1
	INPUT_SLICE
)

type InputDetails struct {
	InputAmount InputAmount
	DataType    string
}

type StatementConfig struct {
	Name         string
	Statement    string
	InputFields  map[string]*InputDetails
	OutputFields []string
	ReturnAmount QueryReturnAmount
}

type Config map[string]StatementConfig

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

	stmts := strings.Split(fileContents, ";")
	nonEmptyStmts := []string{}
	for _, stmt := range stmts {
		if stmt != "" {
			nonEmptyStmts = append(nonEmptyStmts, stmt)
		}
	}

	return nonEmptyStmts, nil
}

func PrepareConfig(ctx context.Context, ddlConfig ddl.Config) (Config, error) {
	config, err := classifyQueryAndPrepare(ctx, rawStatements, ddlConfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}
