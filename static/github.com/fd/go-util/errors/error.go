package errors

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
)

var StackLimit = 3
var StackContext = 3

type Error struct {
	err     error
	message string
	context []string
	stack   Stack
	fatal   bool
}

func IsFatal(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*Error); ok {
		return e.fatal
	}

	return true
}

// Capture any panic() in f() and return it as an error.
func Guard(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = NewFromPanic(r)
		}
	}()

	err = f()
	return
}

// Make a new Error from a panic()
func NewFromPanic(r interface{}) *Error {
	if r == nil {
		return nil
	}

	if err, ok := r.(*Error); ok {
		return Annotate(err, "panic: %s", err.message)
	}

	if err, ok := r.(error); ok {
		return Annotate(err, "panic: %s", err)
	}

	if msg, ok := r.(string); ok {
		return New("panic: %s", msg)
	}

	return New("panic: %+v", r)
}

// Make a new Error
func New(message string, args ...interface{}) *Error {
	return &Error{
		message: fmt.Sprintf(message, args...),
		stack:   CaptureStack().Skip(2),
		fatal:   true,
	}
}

// Wrap err in an new Error
func Annotate(err error, message string, args ...interface{}) *Error {
	if l, ok := err.(List); ok {
		err = l.Normalize()
	}

	if err == nil {
		return nil
	}

	return &Error{
		err:     err,
		message: fmt.Sprintf(message, args...),
		stack:   CaptureStack().Skip(2),
		fatal:   true,
	}
}

func (e *Error) SetFatal(flag bool) {
	if e == nil {
		return
	}

	e.fatal = flag
}

// Add some context to an error
func (e *Error) AddContext(format string, args ...interface{}) {
	if e == nil {
		return
	}

	e.context = append(
		e.context,
		fmt.Sprintf(format, args...),
	)
}

func (e *Error) Error() string {
	if e == nil {
		return "(no error)"
	}

	var (
		buf            bytes.Buffer
		basic_child    bool
		extended_child bool
	)

	switch e.err.(type) {
	case *Error, List:
		extended_child = true
	default:
		if e.err != nil {
			basic_child = true
		}
	}

	fmt.Fprintf(&buf, "error: %s\n", e.message)

	if basic_child {
		fmt.Fprintf(&buf, "  message: %s\n", e.err)
	}

	if len(e.context) > 0 {
		fmt.Fprintln(&buf, "  context:")

		sort.Strings(e.context)

		w := tabwriter.NewWriter(&buf, 0, 8, 1, ' ', 0)

		for _, pair := range e.context {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 1 {
				fmt.Fprintf(w, "    %s\t\n", parts[0])
			} else {
				fmt.Fprintf(w, "    %s\t= %s\n", parts[0], parts[1])
			}
		}

		w.Flush()
	}

	if len(e.stack) > 0 {
		stack := e.stack.Limit(StackLimit).String()
		stack = strings.Replace(stack, "\n", "\n    ", -1)
		fmt.Fprintf(&buf, "  location:\n    %s\n", stack)
	}

	if extended_child {
		err := e.err.Error()
		err = strings.TrimSpace(err)
		err = strings.Replace(err, "\n", "\n  ", -1)
		fmt.Fprintf(&buf, "  %s\n", err)
	}

	return buf.String()
}
