package errors

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
	"text/tabwriter"
)

type Stack []StackFrame

type StackFrame struct {
	PC       uintptr
	Package  string
	Function string
	Filepath string
	Filename string
	Line     int
	InApp    bool

	HasContext  bool
	PreContext  []string
	ContextLine string
	PostContext []string
}

func (s Stack) Skip(n int) Stack {
	if len(s) <= n {
		return nil
	}

	return s[n:]
}

func (s Stack) Limit(n int) Stack {
	if len(s) <= n {
		return s
	}

	return s[:n]
}

func (s Stack) String() string {
	var (
		buf bytes.Buffer
	)

	s.write_to(&buf)
	str := buf.String()
	str = strings.TrimSuffix(str, "\n")

	return str
}

func (s Stack) write_to(w *bytes.Buffer) {
	for _, frame := range s {
		frame.write_to(w)
	}
}

func (f *StackFrame) write_to(w *bytes.Buffer) {
	fmt.Fprintf(w, "%s/%s:%d %s.%s() (0x%x)\n", f.Package, f.Filename, f.Line, path.Base(f.Package), f.Function, f.PC)

	if !f.HasContext {
		return
	}

	var (
		tw = tabwriter.NewWriter(w, 1, 1, 1, ' ', 0)
		s  = f.Line - len(f.PreContext)
	)

	for i, line := range f.PreContext {
		fmt.Fprintf(tw, " \t%d\t%s\n", s+i, strings.Replace(line, "\t", "    ", -1))
	}

	fmt.Fprintf(tw, ">\t%d\t%s\n", f.Line, strings.Replace(f.ContextLine, "\t", "    ", -1))

	for i, line := range f.PostContext {
		fmt.Fprintf(tw, " \t%d\t%s\n", f.Line+1+i, strings.Replace(line, "\t", "    ", -1))
	}

	tw.Flush()
}

func CaptureStack() Stack {
	var (
		stack    []StackFrame
		lines    [][]byte
		lastFile string
	)

	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		frame := StackFrame{
			PC:       pc,
			Filepath: file,
			Line:     line,
			InApp:    true,
		}

		function(pc, &frame)
		filename(&frame)

		// load file
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err == nil {
				lines = bytes.Split(data, []byte{'\n'})
				lastFile = file
			} else {
				lastFile = file
				lines = nil
			}
		}

		if lines != nil {
			line-- // in stack trace, lines are 1-indexed but our array is 0-indexed
			frame.HasContext = true
			context_line(lines, line, &frame)
			pre_context(lines, line, &frame)
			post_context(lines, line, &frame)
		}

		stack = append(stack, frame)
	}

	return Stack(stack)
}

func pre_context(lines [][]byte, n int, frame *StackFrame) {
	i := n - StackContext
	if i < 0 {
		i = 0
	}
	if i == n {
		return
	}

	frame.PreContext = make([]string, 0, StackContext)
	for ; i < n; i++ {
		l := source(lines, i)
		if l == nil {
			continue
		}

		frame.PreContext = append(frame.PreContext, string(l))
	}
}

func post_context(lines [][]byte, n int, frame *StackFrame) {
	i := n + StackContext + 1
	if i > len(lines) {
		i = len(lines)
	}
	if i == n {
		return
	}
	n++

	frame.PostContext = make([]string, 0, StackContext)
	for ; n < i; n++ {
		l := source(lines, n)
		if l == nil {
			continue
		}

		frame.PostContext = append(frame.PostContext, string(l))
	}
}

func context_line(lines [][]byte, n int, frame *StackFrame) {
	frame.ContextLine = string(source(lines, n))
}

func source(lines [][]byte, n int) []byte {
	if n < 0 || n >= len(lines) {
		return nil
	}
	return lines[n]
}

func function(pc uintptr, frame *StackFrame) {
	f := runtime.FuncForPC(pc)
	if f == nil {
		return
	}

	fname := f.Name()
	slash_idx := strings.LastIndex(fname, "/")
	dot_idx := strings.Index(fname[slash_idx+1:], ".")
	pkgpath, name := fname[:slash_idx+1+dot_idx], fname[slash_idx+dot_idx+2:]

	frame.Package = pkgpath
	frame.Function = name
}

func filename(frame *StackFrame) {
	frame.Filename = path.Join(frame.Package, path.Base(frame.Filepath))
}
