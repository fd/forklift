package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fd/forklift/static/github.com/knieriem/markdown"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

var (
	mmparser = markdown.NewParser(nil)
)

func (m *Manual) Open() {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("groffer", "--tty")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	m.format_man(w)
	w.Close()
	r.Close()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

func (m *Manual) format_man(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	buf := bufio.NewWriter(w)
	defer buf.Flush()

	fmt.Fprintf(buf, ".TH command 1\n")
	m.name(buf)
	m.synopsis(buf)
	m.flags(buf)
	m.args(buf)
	m.vars(buf)
	m.commands(buf)

	for _, s := range m.sections {
		s.Man(buf)
	}
}

func (m *Manual) name(w *bufio.Writer) {
	s := section_t{
		Header: "NAME",
		Body:   m.summary,
	}
	s.Man(w)
}

func (m *Manual) synopsis(w *bufio.Writer) {
	s := section_t{
		Header: "SYNOPSIS",
		Body:   m.usage,
	}
	s.Man(w)
}

func (m *Manual) flags(w *bufio.Writer) {
	var (
		buf   bytes.Buffer
		parts []string
		exec  = m.exec
	)

	mmf := markdown.ToGroffMM(&buf)

	for exec != nil {
		for _, f := range exec.Flags {
			s, p := m.options[f.Name]
			if !p {
				continue
			}

			buf.Reset()
			mmparser.Markdown(bytes.NewReader([]byte(s.Body)), mmf)
			body := buf.String()
			body = strings.TrimPrefix(body, ".P\n")

			part := fmt.Sprintf(".TP\n\\fB%s\\fP\n%s", f.Tag.Get("flag"), body)
			parts = append(parts, part)
		}
		exec = exec.ParentExec
	}

	if len(parts) > 0 {
		fmt.Fprintf(w, ".SH OPTIONS\n")
		sort.Strings(parts)
		for _, part := range parts {
			fmt.Fprint(w, part)
		}
		fmt.Fprintf(w, ".P\n")
	}
}

func (m *Manual) vars(w *bufio.Writer) {
	var (
		buf   bytes.Buffer
		parts []string
		exec  = m.exec
	)

	mmf := markdown.ToGroffMM(&buf)

	for exec != nil {
		for _, f := range exec.Variables {
			s, p := m.options[f.Name]
			if !p {
				continue
			}

			if flag := f.Tag.Get("flag"); flag != "" {
				s.Body = fmt.Sprintf("See %s", flag)
			}

			buf.Reset()
			mmparser.Markdown(bytes.NewReader([]byte(s.Body)), mmf)
			body := buf.String()
			body = strings.TrimPrefix(body, ".P\n")

			part := fmt.Sprintf(".TP\n\\fB%s\\fP\n%s", f.Tag.Get("env"), body)
			parts = append(parts, part)
		}
		exec = exec.ParentExec
	}

	if len(parts) > 0 {
		fmt.Fprintf(w, ".SH \"ENVIRONMENT VARIABLES\"\n")
		sort.Strings(parts)
		for _, part := range parts {
			fmt.Fprint(w, part)
		}
		fmt.Fprintf(w, ".P\n")
	}
}

func (m *Manual) args(w *bufio.Writer) {
	var (
		buf   bytes.Buffer
		parts []string
	)

	mmf := markdown.ToGroffMM(&buf)

	for _, f := range m.exec.Args {
		s, p := m.options[f.Name]
		if !p {
			continue
		}

		buf.Reset()
		mmparser.Markdown(bytes.NewReader([]byte(s.Body)), mmf)
		body := buf.String()
		body = strings.TrimPrefix(body, ".P\n")

		part := fmt.Sprintf(".TP\n\\fB%s\\fP\n%s", f.Name, body)
		parts = append(parts, part)
	}

	if len(parts) > 0 {
		fmt.Fprintf(w, ".SH ARGUMENTS\n")
		sort.Strings(parts)
		for _, part := range parts {
			fmt.Fprint(w, part)
		}
		fmt.Fprintf(w, ".P\n")
	}
}

func (m *Manual) commands(w *bufio.Writer) {
	var (
		buf   bytes.Buffer
		parts []string
	)

	mmf := markdown.ToGroffMM(&buf)

	for _, subcmd := range m.exec.SubCommands {
		subm := subcmd.Manual()
		if subm == nil || subm.summary == "" {
			continue
		}

		buf.Reset()
		mmparser.Markdown(bytes.NewReader([]byte(subm.summary)), mmf)
		body := buf.String()
		body = strings.TrimPrefix(body, ".P\n")

		part := fmt.Sprintf(".TP\n\\fB%s\\fP\n%s", strings.Join(subcmd.Names, ", "), body)
		parts = append(parts, part)
	}

	if len(parts) > 0 {
		fmt.Fprintf(w, ".SH COMMANDS\n")
		sort.Strings(parts)
		for _, part := range parts {
			fmt.Fprint(w, part)
		}
		fmt.Fprintf(w, ".P\n")
	}
}

func (s *section_t) Man(w *bufio.Writer) {
	fmt.Fprintf(w, ".SH %s\n", strconv.Quote(s.Header))

	f := markdown.ToGroffMM(w)
	mmparser.Markdown(bytes.NewReader([]byte(s.Body)), f)
	fmt.Fprintf(w, ".P\n")
}
