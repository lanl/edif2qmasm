// This file is part of edif2qmasm.  It abtracts QMASM code to make it
// easier to work with but still convertible to strings in the end.

package main

import (
	"fmt"
	"sort"
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
	}
	return fmt.Sprintf("%s = %s  # %s\n", c.Var[0], c.Var[1], c.Comment)
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
		lines = append(lines, "  "+ln.String())
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
	}
	return fmt.Sprintf("%s := %s  # %s\n", p.Var, b2s[p.Value], p.Comment)
}

// A QmasmRename indicates that one list of variables should be renamed to
// another list of variables.
type QmasmRename struct {
	Before []string // Old names
	After  []string // New names
}

// String outputs a QmasmRename as a line of QMASM code, including a
// training newline.  In practice, these will be batched and performed
// as a single rename to support swapping.
func (r QmasmRename) String() string {
	lhs := fmt.Sprintf("%v", r.Before)
	lhs = lhs[1 : len(lhs)-1]
	rhs := fmt.Sprintf("%v", r.After)
	rhs = rhs[1 : len(rhs)-1]
	return fmt.Sprintf("%s -> %s\n", lhs, rhs)
}

// QmasmCodeList is a slice of QmasmCode lines.
type QmasmCodeList []QmasmCode

// sortPriority maps a QmasmCode to an integer representing its sort priority.
// Only QmasmCode types we expect to find within a macro definition are
// included here.  Anything else aborts with an internal error.
func sortPriority(q QmasmCode) int {
	switch q.(type) {
	case QmasmMacroUse:
		return 0
	case QmasmChain:
		return 1
	case QmasmPin:
		return 2
	default:
		notify.Fatalf("Internal error assigning a priority to %#v", q)
	}
	return 100 // Will never get here
}

// Len returns the length of a QmasmCodeList.  It is used to implement
// sort.Interface.
func (qcl QmasmCodeList) Len() int { return len(qcl) }

// Less says if one QmasmCode line is less than another.  It is used to
// implement sort.Interface.
func (qcl QmasmCodeList) Less(i, j int) bool {
	// Sort first by priority then by string representation.
	ip := sortPriority(qcl[i])
	jp := sortPriority(qcl[j])
	switch {
	case ip < jp:
		return true
	case ip > jp:
		return false
	default:
		return qcl[i].String() < qcl[j].String()
	}
}

// Swap swaps two elements of a QmasmCodeList.  It is used to
// implement sort.Interface.
func (qcl QmasmCodeList) Swap(i, j int) { qcl[i], qcl[j] = qcl[j], qcl[i] }

// SortAndMerge sorts a QmasmCodeList and merges all QmasmRename elements into
// a single QmasmRename.  It returns a new QmasmCodeList.
func (qcl QmasmCodeList) SortAndMerge() QmasmCodeList {
	// Handle the trivial case first.
	if len(qcl) == 0 {
		return qcl
	}

	// Retain all elements except QmasmRename elements.
	before := make([]string, 0, 4) // Old symbol names (4 is arbitrary)
	after := make([]string, 0, 4)  // New symbol names (4 is arbitrary)
	qcl2 := make(QmasmCodeList, 0, len(qcl))
	for _, e := range qcl {
		qr, isQR := e.(QmasmRename)
		if !isQR {
			qcl2 = append(qcl2, e)
			continue
		}
		before = append(before, qr.Before...)
		after = append(after, qr.After...)
	}

	// Sort the list then append a single, new QmasmRename element.
	sort.Sort(qcl2)
	if len(before) > 0 {
		qcl2 = append(qcl2, QmasmRename{
			Before: before,
			After:  after,
		})
	}
	return qcl2
}
