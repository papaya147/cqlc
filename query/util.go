package query

import "github.com/papaya147/go-cassandra-codegen/util"

func getMatch(regex, stmt string) string {
	match, err := util.GetFirstMatch(regex, stmt)
	if err != nil {
		errorList.Add(err)
	}

	return match
}

func getAllMatches(regex, stmt string) []string {
	matches, err := util.GetAllMatches(regex, stmt)
	if err != nil {
		errorList.Add(err)
	}

	return matches
}
