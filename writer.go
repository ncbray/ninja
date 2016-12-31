package ninja

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

const indent = "  "

var specialPathChars = regexp.MustCompile("([$: ])")

func escapePath(path string) string {
	return specialPathChars.ReplaceAllString(path, "$$$1")
}

func escapePaths(paths []string) []string {
	out := make([]string, len(paths))
	for i, path := range paths {
		out[i] = escapePath(path)
	}
	return out
}

type RuleOptions struct {
	Description string
	// "gcc" or "msvc"
	Deps           string
	DepFile        string
	Generator      bool
	Pool           string
	Restat         bool
	RSPFile        string
	RSPFileContent string
}

type BuildOptions struct {
	Inputs          []string
	ImplicitInputs  []string
	OrderOnlyInputs []string
	ImplicitOutputs []string
}

type NinjaWriter struct {
	out io.Writer
}

func (w *NinjaWriter) line(text string, indent_count int) {
	for i := 0; i < indent_count; i++ {
		io.WriteString(w.out, indent)
	}
	io.WriteString(w.out, text)
	io.WriteString(w.out, "\n")
}

func (w *NinjaWriter) variable(name string, value string) {
	if value != "" {
		w.line(fmt.Sprintf("%s = %s", name, value), 1)
	}
}

func (w *NinjaWriter) Rule(name string, cmd string, options RuleOptions) {
	w.line(fmt.Sprintf("rule %s", name), 0)
	w.variable("command", cmd)

	w.variable("description", options.Description)
	w.variable("deps", options.Deps)
	w.variable("depfile", options.DepFile)
	if options.Generator {
		w.variable("generator", "1")
	}
	w.variable("pool", options.Pool)
	if options.Restat {
		w.variable("restat", "1")
	}

	w.variable("rspfile", options.RSPFile)
	w.variable("rspfile_content", options.RSPFileContent)
}

func appendEscapedPaths(existing []string, prefix string, additional []string) []string {
	if len(additional) > 0 {
		existing = append(append(existing, prefix), escapePaths(additional)...)
	}
	return existing
}

func (w *NinjaWriter) Build(outputs []string, rule string, options BuildOptions) {
	outputs = escapePaths(outputs)
	outputs = appendEscapedPaths(outputs, "|", options.ImplicitOutputs)

	inputs := escapePaths(options.Inputs)
	inputs = appendEscapedPaths(inputs, "|", options.ImplicitInputs)
	inputs = appendEscapedPaths(inputs, "||", options.OrderOnlyInputs)

	w.line(fmt.Sprintf(
		"build %s: %s",
		strings.Join(outputs, " "),
		strings.Join(append([]string{rule}, inputs...), " "),
	), 0)

	// TODO variable overrides.
}

func MakeNinjaWriter(out io.Writer) *NinjaWriter {
	return &NinjaWriter{
		out: out,
	}
}
