package cli

import (
	"fmt"
	"strings"
	"unicode"
)

func (exec *executable_t) Manual() *Manual {
	if exec.manual.Type == nil {
		return nil
	}

	if exec.parsed_manual != nil {
		return exec.parsed_manual
	}

	m := &Manual{}
	m.parse(exec)
	exec.parsed_manual = m

	return m
}

func (m *Manual) parse(exec *executable_t) {
	var (
		source = string(exec.manual.Tag)
		indent = determine_indent(source)
		lines  = strings.Split(source, "\n")
	)

	m.exec = exec
	m.options = make(map[string]section_t, 20)

	if p := exec.ParentExec; p != nil {
		if pm := p.Manual(); pm != nil {
			for k, o := range pm.options {
				m.options[k] = o
			}
		}
	}

	remove_indent(lines, indent)
	m.parse_sections(lines)
}

func (m *Manual) parse_sections(lines []string) {
	var (
		in_section   bool
		section_name string
		section_body string
	)

	for _, line := range lines {
		name, body, empty := parse_line(line)

		if empty {
			if in_section {
				section_body += "\n"
			}
			continue
		}

		if name != "" {
			if in_section {
				m.parse_section(section_name, section_body)
				section_name = ""
				section_body = ""
			}

			in_section = true
			section_name = name
		}

		section_body += body + "\n"
	}

	if in_section {
		m.parse_section(section_name, section_body)
	}
}

func (m *Manual) parse_section(name, body string) {
	{
		var (
			indent = determine_indent(body)
			lines  = strings.Split(body, "\n")
		)

		remove_indent(lines, indent)
		body = strings.TrimSpace(strings.Join(lines, "\n"))
	}

	switch name {

	case "Usage":
		m.usage = body

	case "Summary":
		m.summary = body

	default:

		if strings.HasPrefix(name, ".") {
			p := section_t{Header: name, Body: body}
			m.options[name[1:]] = p

		} else {
			p := section_t{Header: name, Body: body}
			m.sections = append(m.sections, p)
		}

	}
}

func parse_line(line string) (section, body string, empty bool) {
	if len(line) == 0 || is_space_only(line) {
		return "", "", true
	}

	if unicode.IsSpace([]rune(line)[0]) {
		return "", line, false
	}

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		panic(fmt.Sprintf("Invalid indentation for line: `%s`", line))
	}

	return strings.TrimSpace(parts[0]), strings.TrimLeft(parts[1], " \t"), false
}

func remove_indent(lines []string, n int) {
	for i, line := range lines {
		lines[i] = skip_at_most_n_spaces(line, n)
	}
}

func determine_indent(source string) int {
	var (
		indent int
	)

	for _, c := range source {
		if c == '\n' {
			indent = 0
		} else if unicode.IsSpace(c) {
			indent += 1
		} else {
			break
		}
	}

	return indent
}

func skip_at_most_n_spaces(line string, n int) string {
	var (
		prefix string
		suffix string
	)

	if len(line) < n {
		prefix = line
	} else {
		prefix = line[:n]
		suffix = line[n:]
	}

	if is_space_only(prefix) {
		return suffix
	}

	panic(fmt.Sprintf("Invalid indentation for line: `%s`", line))
}

func is_space_only(s string) bool {
	for _, c := range s {
		if !unicode.IsSpace(c) {
			return false
		}
	}
	return true
}
