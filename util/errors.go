package util

type ErrorList []error

func NewErrorList() ErrorList {
	return ErrorList{}
}

func (e *ErrorList) Error() string {
	if e.IsEmpty() {
		return ""
	}

	out := ""
	for _, err := range *e {
		out += ", " + err.Error()
	}
	return out[2:]
}

func (e *ErrorList) Add(err error) {
	*e = append(*e, err)
}

func (e *ErrorList) IsEmpty() bool {
	return len(*e) == 0
}
