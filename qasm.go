// This file is part of edif2qasm.  It abtracts QASM code to make it
// easier to work with but still convertible to strings in the end.

package main

import (
	"fmt"
	"strings"
)

// QasmCode is anything defined in this file.  At a minimum, it must
// be convertible to a string.
type QasmCode interface {
	String() string
}

// A QasmVar is a named QASM variable.
type QasmVar string

// A QasmChain indicates that two variables should be assigned the same value.
type QasmChain struct {
	Q1, Q2  QasmVar // Variables to equate
	Comment string  // Optional comment
}

// String outputs a QasmChain as a line of QASM code, including a training
// newline.
func (c QasmChain) String() string {
	if c.Comment == "" {
		return fmt.Sprintf("%s = %s\n", c.Q1, c.Q2)
	} else {
		return fmt.Sprintf("%s = %s  # %s\n", c.Q1, c.Q2, c.Comment)
	}
}

// A QasmMacroDef represents a QASM macro definition.
type QasmMacroDef struct {
	Name    string     // Macro name
	Body    []QasmCode // Macro body
	Comment string     // Optional comment
}

// String outputs a QASM macro definition.
func (m QasmMacroDef) String() string {
	lines := make([]string, 0, 4)
	if m.Comment != "" {
		lines = append(lines, "# "+m.Comment+"\n")
	}
	lines = append(lines, "!begin_macro "+m.Name+"\n")
	for _, ln := range m.Body {
		lines = append(lines, ln.String())
	}
	lines = append(lines, "!end_macro "+m.Name+"\n")
	return strings.Join(lines, "")
}

// A QasmMacroUse instantiates a QASM macro.
type QasmMacroUse struct {
	MacroName string // Name of the macro to instantiate
	UseName   string // Name of the instantiation
	Comment   string // Optional comment
}

// String outputs a QASM macro use.
func (u QasmMacroUse) String() string {
	str := ""
	if u.Comment != "" {
		str += "# " + u.Comment + "\n"
	}
	str += "!use_macro " + u.MacroName + " " + u.UseName
	return str
}

// A QasmComment is a QASM comment with no associated code.
type QasmComment struct {
	Comment string // The comment itself
}

// String outputs a QASM comment with a trailing newline.
func (c QasmComment) String() string {
	return "# " + c.Comment + "\n"
}

// A QasmBlank is a no-op, output as a blank line for aesthetic purposes.
type QasmBlank struct{}

// String outputs a QASM no-op as a single newline.
func (b QasmBlank) String() string {
	return "\n"
}
