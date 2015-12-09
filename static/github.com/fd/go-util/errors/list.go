package errors

import (
	"strings"
)

type List []error

func (l List) Normalize() error {
	if l.HasErrors() {
		return l
	}

	return nil
}

func (l List) HasErrors() bool {
	return len(l) > 0
}

func (l *List) Add(err error) {
	if err == nil {
		return
	}

	if list, ok := err.(List); ok {
		for _, err := range list {
			if err != nil {
				*l = append(*l, err)
			}
		}
		return
	}

	*l = append(*l, err)
}

func (l List) Error() string {
	s := make([]string, len(l))

	for i, err := range l {
		s[i] = err.Error()
	}

	return strings.Join(s, "\n")
}
