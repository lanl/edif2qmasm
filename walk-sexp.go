// This file is part of edif2qasm.  It provides functions for walking an EDIF
// s-expression, converting it to QASM format.

package main

// ConvertMetadata converts top-level metadata to QASM.
func ConvertMetadata(s EdifSExp) []QasmCode {
	hdr := make([]QasmCode, 0, 1)
	el := AsList(s, 1, "edif")
	hdr = append(hdr, QasmComment{
		Comment: "Module " + string(AsSymbol(el[1]))},
	)
	cmts := el.SublistsByName("comment")
	for _, c := range cmts {
		hdr = append(hdr, QasmComment{
			Comment: string(AsString(c[1])),
		})
	}
	hdr = append(hdr, QasmComment{
		Comment: "Converted to QASM by edif2qasm",
	})
	return hdr
}

// CanonicalizeCellName converts names of the form "$_AND_" to the form "AND".
func CanonicalizeCellName(s EdifString) EdifString {
	n := len(s)
	if n < 4 {
		return s
	}
	if s[0] == '$' && s[1] == '_' && s[n-1] == '_' {
		return s[2 : n-1]
	}
	return s
}

// ProcessExternalLib examines external-library imports and returns a map from
// cell ID to cell name.
func ProcessExternalLib(e EdifList) map[EdifSymbol]EdifString {
	idToName := make(map[EdifSymbol]EdifString, 8)
	for _, cell := range e.SublistsByName("cell") {
		cnm := cell[1]
		if cnm.Type() != List {
			continue // Symbols don't need to be mapped.
		}
		rnm := AsList(cnm, 2, "rename")
		idToName[AsSymbol(rnm[1])] = CanonicalizeCellName(AsString(rnm[2]))
	}
	return idToName
}

// ConvertInstance converts an instantiated cell to a QASM macro instantiation.
// This function returns a slice rather than an individual QasmMacroUse because
// it may need to return an empty slice (in the case of VCC and GND).
func ConvertInstance(inst EdifList, i2n map[EdifSymbol]EdifString) []QasmCode {
	// Extract the instantiation name.
	code := make([]QasmCode, 0, 1)
	var instSym EdifSymbol // Instantiated macro name
	var comment string     // Comment describing the instantiation
	switch inst[1].Type() {
	case Symbol:
		instSym = AsSymbol(inst[1])
		if instSym == "GND" || instSym == "VCC" {
			return nil // GND and VCC are treated specially.
		}

	case List:
		ren := AsList(inst[1], 3, "rename")
		instSym = AsSymbol(ren[1])
		comment = string(AsString(ren[2]))
	}

	// Extract the macro name.
	var macroName EdifSymbol
	cRef := inst.NestedSublistsByName([]EdifSymbol{
		"viewRef",
		"cellRef",
	})[0]
	macroName = AsSymbol(cRef[1])
	if cName, ok := i2n[macroName]; ok {
		macroName = EdifSymbol(cName) // Renamed cell (e.g., "id0123" --> "AND")
	}

	// Construct and return a macro instantiation.
	code = append(code, QasmMacroUse{
		MacroName: string(macroName),
		UseName:   "$" + string(instSym),
		Comment:   comment,
	})

	return code
}

// ConvertNet converts an EDIF net to a QASM chain ("=").
func ConvertNet(net EdifList) QasmChain {
	// Determine the name of each port.
	portName := make([]string, 0, 2)
	for _, pRef := range net.NestedSublistsByName([]EdifSymbol{
		"joined",
		"portRef",
	}) {
		switch len(pRef) {
		case 2:
			// Symbol is defined by the current macro.
			portName = append(portName, string(AsSymbol(pRef[1])))

		case 3:
			// Symbol is defined by an instantiated macro.
			pName := string(AsSymbol(pRef[1]))
			instRef := AsList(pRef[2], 2, "instanceRef")
			pName = "$" + string(AsSymbol(instRef[1])) + "." + pName
			portName = append(portName, pName)

		default:
			notify.Fatalf("Expected 2 or 3 elements in a portRef; say %v", pRef)
		}
	}
	if len(portName) != 2 {
		notify.Fatalf("Expected a net to contain exactly two portRefs; saw %v", net)
	}

	// Treat a renamed net as a comment.
	comment := ""
	if net[1].Type() == List {
		ren := AsList(net[1], 3, "rename")
		comment = string(AsString(ren[2]))
	}

	// Return a QASM chain.
	return QasmChain{
		Var:     [2]string{portName[0], portName[1]},
		Comment: comment,
	}
}

// ConvertCell converts a user-defined cell to a QASM macro definition.
func ConvertCell(cell EdifList, i2n map[EdifSymbol]EdifString) QasmMacroDef {
	// Ensure the cell looks at least a little like what we expect.
	if len(cell) < 3 {
		notify.Fatalf("Cell %v contains too few components", cell)
	}

	// Instantiate all the other cells used by the current cell.
	code := make([]QasmCode, 0, 32)
	for _, inst := range cell.NestedSublistsByName([]EdifSymbol{
		"view",
		"contents",
		"instance",
	}) {
		code = append(code, ConvertInstance(inst, i2n)...)
	}

	// Instantiate all the nets used by the current cell.
	for _, net := range cell.NestedSublistsByName([]EdifSymbol{
		"view",
		"contents",
		"net",
	}) {
		code = append(code, ConvertNet(net))
	}

	// Wrap the code in a QASM macro definition and return it.
	return QasmMacroDef{
		Name: string(AsSymbol(cell[1])),
		Body: code,
	}
}

// ConvertDesign converts a design to a QASM macro instantiation.
func ConvertDesign(des EdifList) QasmMacroUse {
	if len(des) != 3 {
		notify.Fatalf("Expected a design to contain exactly 3 elements but saw %v", des)
	}
	cRef := AsList(des[2], 3, "cellRef")
	return QasmMacroUse{
		MacroName: "$" + string(AsSymbol(cRef[1])),
		UseName:   string(AsSymbol(des[1])),
	}
}

// ConvertLibrary converts a user-defined cell library to QASM macro
// definitions.
func ConvertLibrary(lib EdifList, i2n map[EdifSymbol]EdifString) []QasmCode {
	// Iterate over each cell.
	code := make([]QasmCode, 0, 32)
	for _, cell := range lib.SublistsByName("cell") {
		code = append(code, ConvertCell(cell, i2n))
	}
	return code
}

// ConvertEdifToQasm takes an EDIF s-expression and returns a list of QASM
// statements.
func ConvertEdifToQasm(s EdifSExp) []QasmCode {
	// Produce a QASM header block.
	code := make([]QasmCode, 0, 128)
	code = append(code, ConvertMetadata(s)...)
	code = append(code, QasmBlank{})
	code = append(code, QasmInclude{File: "stdcell"})

	// Generate a mapping from cell ID to cell name.
	idToName := make(map[EdifSymbol]EdifString, 8)
	slst := AsList(s, 2, "edif")
	for _, ext := range slst.SublistsByName("external") {
		for id, nm := range ProcessExternalLib(ext) {
			idToName[id] = nm
		}
	}

	// Convert each user-defined library in turn.
	for _, lib := range slst.SublistsByName("library") {
		code = append(code, QasmBlank{})
		code = append(code, ConvertLibrary(lib, idToName)...)
	}

	// Convert each design in turn.
	for _, des := range slst.SublistsByName("design") {
		code = append(code, QasmBlank{})
		code = append(code, ConvertDesign(des))
	}

	return code
}
