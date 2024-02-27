package dml

import (
	"errors"
	"fmt"
	"strings"

	"github.com/papaya147/cqlc/ddl"
	"github.com/papaya147/cqlc/util"
)

const captureUpdateKeyspaceAndTable = `(?is)UPDATE\s+([A-Za-z_0-9.]+)`
const captureUpdateReturnValues = `(?is)RETURNING\s+(.*)`
const captureUpdateSingleInputValues = `(?i)([A-Za-z_0-9()]+)\s*[=><]+\s*\?`
const captureUpdateSliceInputValues = `(?i)([A-Za-z_0-9()]+)\s+IN\s*\(%s`

func loadUpdateStatement(stmt string, ddlConfig ddl.Config, out *StatementConfig) error {
	list := util.NewErrorList()

	keyspaceAndTableNamesBlock, err := util.GetFirstMatch(captureUpdateKeyspaceAndTable, stmt)
	if err != nil {
		list.Add(errors.New("unable to extract keyspace and name from query"))
	}

	keyspaceAndTableName := strings.Split(keyspaceAndTableNamesBlock, ".")
	if len(keyspaceAndTableName) != 2 {
		list.Add(errors.New("unable to extract keyspace and name from query"))
		return list
	}

	keyspace := keyspaceAndTableName[0]
	tableName := keyspaceAndTableName[1]

	// check if table exists in the keyspace
	if tableConfig, ok := ddlConfig[tableName]; ok {
		if tableConfig.Keyspace != keyspace {
			list.Add(fmt.Errorf("table %s does not exist in keyspace %s", tableName, keyspace))
		}
	} else {
		list.Add(fmt.Errorf("table %s does not exist", tableName))
	}

	inputValues := make(map[string]*InputDetails)

	singleInputValues, _ := util.GetAllMatches(captureUpdateSingleInputValues, stmt)
	for _, value := range singleInputValues {
		val := strings.TrimSpace(strings.ToLower(value))
		if inputValues[val] != nil {
			j := 1
			newFieldName := fmt.Sprintf("%s%d", val, j)
			for inputValues[newFieldName] != nil {
				j++
			}
			inputValues[newFieldName] = &InputDetails{
				InputAmount: INPUT_ONE,
				DataType:    userOptions.TypeMappings[ddlConfig[tableName].Fields[val]],
			}
		} else {
			inputValues[val] = &InputDetails{
				InputAmount: INPUT_ONE,
				DataType:    userOptions.TypeMappings[ddlConfig[tableName].Fields[val]],
			}
		}
	}

	sliceInputValues, _ := util.GetAllMatches(captureUpdateSliceInputValues, stmt)
	for _, value := range sliceInputValues {
		val := strings.TrimSpace(strings.ToLower(value))
		inputValues[val] = &InputDetails{
			InputAmount: INPUT_SLICE,
			DataType:    userOptions.TypeMappings[ddlConfig[tableName].Fields[val]],
		}
	}

	outputValuesBlock, err := util.GetAllMatches(captureUpdateReturnValues, stmt)
	var outputValues []string
	if err == nil {
		outputValues = strings.Split(outputValuesBlock[0], ",")
		for i, value := range outputValues {
			outputValues[i] = strings.TrimSpace(value)
		}
		if len(outputValues) == 0 {
			list.Add(errors.New("unable to extract output values from query"))
		}
		if len(outputValues) == 1 && strings.TrimSpace(outputValues[0]) == "*" {
			tableFields := make([]string, len(ddlConfig[tableName].Fields))
			i := 0
			for field := range ddlConfig[tableName].Fields {
				tableFields[i] = field
				i++
			}
			outputValues = tableFields
		}
	}

	if !list.IsEmpty() {
		return list
	}

	config := StatementConfig{
		Statement:    stmt,
		InputFields:  inputValues,
		OutputFields: outputValues,
	}

	*out = config

	return nil
}
