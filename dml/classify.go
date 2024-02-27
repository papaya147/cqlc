package dml

import (
	"context"
	"fmt"
	"strings"

	"github.com/papaya147/cqlc/ddl"
	"github.com/papaya147/cqlc/util"
	"github.com/papaya147/parallelize"
)

const captureQueryName = `(?i)-- name: ([A-Za-z_][A-Za-z_0-9]+)`
const captureQueryReturnAmount = `(?i)-- name: [A-Za-z_][A-Za-z_0-9]+ :([a-z]+)`
const captureQueryType = `(?i)-- name: [A-Za-z_][A-Za-z_0-9]+ :[a-z]+\s+([A-Za-z]+)`
const captureQueryStatement = `(?is)-- name: [A-Za-z_][A-Za-z_0-9]+ :[a-z]+\s+(.*)`

func classifyQueryAndPrepare(ctx context.Context, stmts []string, ddlConfig ddl.Config) (Config, error) {
	config := make(Config)
	configChannel := make(chan StatementConfig, len(stmts))
	group := parallelize.NewSyncGroup()

	for _, stmt := range stmts {
		parallelize.AddMethodWithArgs(group, loadQueryConfig, parallelize.MethodWithArgsParams[loadQueryConfigInput]{
			Context: ctx,
			Input: loadQueryConfigInput{
				stmt:       stmt,
				ddlConfig:  ddlConfig,
				outChannel: configChannel,
			},
		})
	}

	if err := group.Run(); err != nil {
		return nil, err
	}

	close(configChannel)
	for stmtConfig := range configChannel {
		config[stmtConfig.Name] = stmtConfig
	}

	return config, nil
}

type loadQueryConfigInput struct {
	stmt       string
	ddlConfig  ddl.Config
	outChannel chan StatementConfig
}

func loadQueryConfig(ctx context.Context, input loadQueryConfigInput) error {
	list := util.NewErrorList()

	queryName, err := util.GetFirstMatch(captureQueryName, input.stmt)
	if err != nil {
		list.Add(fmt.Errorf("query name not present"))
	}
	queryName = strings.TrimSpace(queryName)

	var ra QueryReturnAmount
	returnAmount, err := util.GetFirstMatch(captureQueryReturnAmount, input.stmt)
	if err != nil {
		list.Add(fmt.Errorf("return amount not present on query %s", queryName))
	}
	switch strings.ToLower(strings.TrimSpace(returnAmount)) {
	case "exec":
		ra = RETURN_NONE
	case "one":
		ra = RETURN_ONE
	case "many":
		ra = RETURN_MANY
	default:
		list.Add(fmt.Errorf("unknown return amount %s", returnAmount))
		return list
	}

	queryType, err := util.GetFirstMatch(captureQueryType, input.stmt)
	if err != nil {
		list.Add(fmt.Errorf("unable to extract query type from query %s", queryName))
	}

	queryStatement, err := util.GetFirstMatch(captureQueryStatement, input.stmt)
	if err != nil {
		list.Add(fmt.Errorf("unable to extract query statement from query %s", queryName))
	}

	var stmtConfig StatementConfig
	switch strings.ToLower(strings.TrimSpace(queryType)) {
	case "select":
		if err := loadSelectStatement(queryStatement, input.ddlConfig, &stmtConfig); err != nil {
			list.Add(err)
			return list
		}
	case "update":
		if err := loadUpdateStatement(queryStatement, input.ddlConfig, &stmtConfig); err != nil {
			list.Add(err)
			return list
		}
	default:
		list.Add(fmt.Errorf("unknown query type %s (insert and delete queries are not yet supported)", queryType))
		return list
	}

	input.outChannel <- StatementConfig{
		Name:         queryName,
		Statement:    stmtConfig.Statement,
		InputFields:  stmtConfig.InputFields,
		OutputFields: stmtConfig.OutputFields,
		ReturnAmount: ra,
	}

	if !list.IsEmpty() {
		return list
	}

	return nil
}
