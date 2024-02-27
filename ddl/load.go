package ddl

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

type TableConfig struct {
	Keyspace   string
	Fields     Fields
	PrimaryKey Keys
}

type Config map[string]TableConfig

func Load(ctx context.Context, opts *options.Config) error {
	stmts, err := getDDL(ctx, opts.Cql.SchemaDir)
	if err != nil {
		return err
	}

	userOptions = opts
	rawStatements = stmts

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

	stmts := strings.Split(fileContents, ";")
	nonEmptyStmts := []string{}
	for _, stmt := range stmts {
		if stmt != "" {
			nonEmptyStmts = append(nonEmptyStmts, stmt)
		}
	}

	return nonEmptyStmts, nil
}

func PrepareConfig(ctx context.Context) (Config, error) {
	// loading keyspaces and tables
	loadKeyspaceNames(ctx)
	if err := loadTableNames(ctx); err != nil {
		return nil, err
	}
	if err := loadTableFields(ctx, userOptions.TypeMappings); err != nil {
		return nil, err
	}

	config := make(Config)

	for table, fields := range tableFields {
		config[table] = TableConfig{
			Keyspace:   tableKeyspaceMap[table],
			Fields:     fields,
			PrimaryKey: tableKeys[table],
		}
	}

	return config, nil
}
