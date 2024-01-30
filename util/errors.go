package util

import "errors"

type ErrorList []error

func NewErrorList() ErrorList {
	return ErrorList{}
}

func (e *ErrorList) SerialiseString() string {
	out := ""
	for _, err := range *e {
		out += err.Error() + "\n"
	}
	return out
}

func (e *ErrorList) Add(err error) {
	*e = append(*e, err)
}

func (e *ErrorList) SerialiseError() error {
	out := ""
	for _, err := range *e {
		out += err.Error() + "\n"
	}
	return errors.New(out)
}

func (e *ErrorList) IsEmpty() bool {
	return len(*e) == 0
}
