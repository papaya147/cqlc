package query

type QueryReturnAmount int

const (
	RETURN_UNSUPPORTED QueryReturnAmount = iota
	RETURN_NONE
	RETURN_ONE
	RETURN_MANY
)

type QueryType int

const (
	QUERY_UNSUPPORTED QueryType = iota
	QUERY_SELECT
	QUERY_INSERT
	QUERY_UPDATE
	QUERY_DELETE
)

type ArgType int

const (
	ARG_SINGLE ArgType = iota
	ARG_LIST
)

func getKeyspaceExtractor(typ QueryType) string {
	switch typ {
	case QUERY_SELECT:
		return `from ([a-z0-9_]+)`
	case QUERY_INSERT:
		return `insert into ([a-z0-9_]+)`
	case QUERY_UPDATE:
		return `update ([a-z0-9_]+)`
	case QUERY_DELETE:
		return `from ([a-z0-9_]+)`
	default:
		return ""
	}
}

func getTableExtractor(typ QueryType) string {
	switch typ {
	case QUERY_SELECT:
		return `from [a-z0-9_]+.([a-z0-9_]+)`
	case QUERY_INSERT:
		return `insert into [a-z0-9_]+.([a-z0-9_]+)`
	case QUERY_UPDATE:
		return `update [a-z0-9_]+.([a-z0-9_]+)`
	case QUERY_DELETE:
		return `from [a-z0-9_]+.([a-z0-9_]+)`
	default:
		return ""
	}
}

func getReturnFieldsExtractor(typ QueryType) string {
	switch typ {
	case QUERY_SELECT:
		return `select\s+([a-z0-9_(),\s])*\s+from`
	case QUERY_INSERT:
		return `returning\s+([a-z0-9_(),\s])*`
	case QUERY_UPDATE:
		return `returning\s+([a-z0-9_(),\s])*`
	case QUERY_DELETE:
		return `returning\s+([a-z0-9_(),\s])*`
	default:
		return ""
	}
}
