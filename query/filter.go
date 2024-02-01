package query

import (
	"fmt"
	"strings"

	"github.com/papaya147/go-cassandra-codegen/util"
)

const validQueryPattern = `-- name:\s[a-zA-Z0-9]+\s:[a-z]+\s*[^;]*`

func (q *QueryList) loadQueries(dml ...string) {
	for _, stmt := range dml {
		if !util.CheckMatch(validQueryPattern, stmt) {
			continue
		}

		q.getQuery(stmt)
	}
}

func (q *QueryList) getQuery(stmt string) {
	queryName := getMatch(queryNameFromQueryComment, stmt)
	returnAmount := getReturnAmount(stmt)

	stmt = strings.ToLower(stmt)

	queryType := getQueryType(stmt)

	keyspace := getMatch(getKeyspaceExtractor(queryType), stmt)

	// checking if keyspace is in schema
	if !util.Contains[string](q.schema.Keyspaces, keyspace) {
		errorList.Add(fmt.Errorf("keyspace %s not found in create statements", keyspace))
	}

	table := getMatch(getTableExtractor(queryType), stmt)

	// checking if table is in schema
	if _, ok := q.schema.Tables[table]; !ok {
		errorList.Add(fmt.Errorf("table %s not found in create statements", keyspace))
	}

	returnFields := getReturnFields(stmt, queryType)

	args := getAllMatches(operatorComparisionsFromQuery, stmt)

	// checking if args are in schema
	argM := []argMap{}
	for _, arg := range args {
		field := getMatch(`(^[^\s]+)`, arg)
		if _, ok := q.schema.Tables[table].Fields[field]; !ok {
			errorList.Add(fmt.Errorf("field %s not found in table %s", field, table))
		}

		operator := getMatch(`^[^\s]+\s+([^\s]+)`, arg)
		if operator == "in" {
			argM = append(argM, argMap{
				Argument: field,
				ArgType:  ARG_LIST,
			})
		} else {
			argM = append(argM, argMap{
				Argument: field,
				ArgType:  ARG_SINGLE,
			})
		}
	}

	q.Query[queryName] = queryDetails{
		Keyspace:     keyspace,
		Table:        table,
		Arguments:    argM,
		ReturnFields: returnFields,
		ReturnAmount: returnAmount,
	}
}

const queryNameFromQueryComment = `-- name:\s([a-zA-Z0-9]+)\s:[a-z]+`

const returnAmountFromQueryComment = `-- name:\s[a-zA-Z0-9]+\s:([a-z]+)`

func getReturnAmount(stmt string) QueryReturnAmount {
	match := getMatch(returnAmountFromQueryComment, stmt)

	switch match {
	case "exec":
		return RETURN_NONE
	case "one":
		return RETURN_ONE
	case "many":
		return RETURN_MANY
	default:
		errorList.Add(fmt.Errorf("%s return type not supported", match))
		return RETURN_UNSUPPORTED
	}
}

const queryTypeFromQuery = `-- name:[^:]+:[a-z]+\s*([a-z]+)`

func getQueryType(stmt string) QueryType {
	match := getMatch(queryTypeFromQuery, stmt)

	switch match {
	case "select":
		return QUERY_SELECT
	case "insert":
		return QUERY_INSERT
	case "update":
		return QUERY_UPDATE
	case "delete":
		return QUERY_DELETE
	default:
		errorList.Add(fmt.Errorf("%s query type not supported", match))
		return QUERY_UNSUPPORTED
	}
}

func getReturnFields(stmt string, typ QueryType) []string {
	match := getMatch(getReturnFieldsExtractor(typ), stmt)
	fmt.Println(match)

	fields := strings.Split(strings.TrimSpace(match), ",")

	var result []string
	for _, field := range fields {
		result = append(result, strings.TrimSpace(field))
	}

	return result
}

const operatorComparisionsFromQuery = `([^\s]+\s*[=><]+\s*[^\s]+|[^\s]+\s+in\s+[^\s]+)`
