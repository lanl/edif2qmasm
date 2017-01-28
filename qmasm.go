// This file is part of edif2qmasm.  It abtracts QMASM code to make it
// easier to work with but still convertible to strings in the end.

package main

import (
	"fmt"
	"strings"
)

// QmasmCode is anything defined in this file.  At a minimum, it must
// be convertible to a string.
type QmasmCode interface {
	String() string
}

// A QmasmChain indicates that two variables should be assigned the same value.
type QmasmChain struct {
	Var     [2]string // Variables to equate (no implied order)
	Comment string    // Optional comment
}

// String outputs a QmasmChain as a line of QMASM code, including a training
// newline.
func (c QmasmChain) String() string {
	if c.Comment == "" {
		return fmt.Sprintf("%s = %s\n", c.Var[0], c.Var[1])
	} else {
		return fmt.Sprintf("%s = %s  # %s\n", c.Var[0], c.Var[1], c.Comment)
	}
}

// A QmasmMacroDef represents a QMASM macro definition.
type QmasmMacroDef struct {
	Name    string      // Macro name
	Body    []QmasmCode // Macro body
	Comment string      // Optional comment
}

// String outputs a QMASM macro definition.
func (m QmasmMacroDef) String() string {
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

// A QmasmMacroUse instantiates a QMASM macro.
type QmasmMacroUse struct {
	MacroName string   // Name of the macro to instantiate
	UseNames  []string // Name(s) of the instantiation
	Comment   string   // Optional comment
}

// String outputs a QMASM macro use.
func (u QmasmMacroUse) String() string {
	str := "!use_macro " + u.MacroName + " " + strings.Join(u.UseNames, " ")
	if u.Comment != "" {
		str += "  # " + u.Comment
	}
	str += "\n"
	return str
}

// A QmasmComment is a QMASM comment with no associated code.
type QmasmComment struct {
	Comment string // The comment itself
}

// String outputs a QMASM comment with a trailing newline.
func (c QmasmComment) String() string {
	return "# " + c.Comment + "\n"
}

// A QmasmBlank is a no-op, output as a blank line for aesthetic purposes.
type QmasmBlank struct{}

// String outputs a QMASM no-op as a single newline.
func (b QmasmBlank) String() string {
	return "\n"
}

// A QmasmInclude includes an external QMASM file.
type QmasmInclude struct {
	File    string // Name of file to include
	Comment string // Optional comment
}

// String outputs a QMASM include with a trailing newline.  We assume the
// external file exists in the user's QMASMPATH.
func (i QmasmInclude) String() string {
	return "!include <" + i.File + ">\n"
}

// A QmasmPin indicates that variables should be pinned to either TRUE or FALSE.
type QmasmPin struct {
	Var     string // Variable to pin
	Value   bool   // Value to assign to the variable
	Comment string // Optional comment
}

// String outputs a QmasmPin as a line of QMASM code, including a training
// newline.
func (p QmasmPin) String() string {
	b2s := map[bool]string{
		false: "GND",
		true:  "VCC",
	}
	if p.Comment == "" {
		return fmt.Sprintf("%s := %s\n", p.Var, b2s[p.Value])
	} else {
		return fmt.Sprintf("%s := %s  # %s\n", p.Var, b2s[p.Value], p.Comment)
	}
}
