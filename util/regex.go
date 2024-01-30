package util

import (
	"fmt"
	"regexp"
)

func CheckMatch(regex, str string) bool {
	re := regexp.MustCompile(regex)
	return re.MatchString(str)
}

func GetFirstMatch(regex, str string) (string, error) {
	if !CheckMatch(regex, str) {
		return "", fmt.Errorf("%s not found in %s", regex, str)
	}

	re := regexp.MustCompile(regex)

	matches := re.FindStringSubmatch(str)
	if len(matches) >= 2 {
		return matches[1], nil
	}

	return "", fmt.Errorf("%s not found in %s", regex, str)
}
