// This file is part of edif2qmasm.  It provides functions for walking an EDIF
// s-expression, converting it to QMASM format.

package main

import (
	"fmt"
	"regexp"
)

// isFlipFlop indicates that a given macro name represents a flip-flop.
var isFlipFlop = make(map[EdifSymbol]bool)

// Given a list of the form (<anything> <name>) or (<anything> (rename
// <name> <comment>)), extract and return <name> and <comment>.
func nameAndComment(e EdifSExp) (EdifSymbol, string) {
	var name EdifSymbol
	var comment string
	if e.Type() == List {
		ren := AsList(e, 3, "rename")
		name = AsSymbol(ren[1])
		comment = string(AsString(ren[2]))
	} else {
		name = AsSymbol(e)
	}
	if len(comment) >= 2 && comment[0] == '\\' {
		comment = comment[1:]
	}
	return name, comment
}

// ConvertMetadata converts top-level metadata to QMASM.
func ConvertMetadata(s EdifSExp) []QmasmCode {
	hdr := make([]QmasmCode, 0, 1)
	el := AsList(s, 1, "edif")
	modName, _ := nameAndComment(el[1])
	hdr = append(hdr, QmasmComment{
		Comment: "Module " + string(modName),
	})
	cmts := el.SublistsByName("comment")
	for _, c := range cmts {
		hdr = append(hdr, QmasmComment{
			Comment: string(AsString(c[1])),
		})
	}
	hdr = append(hdr, QmasmComment{
		Comment: "Converted to QMASM by edif2qmasm",
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

// ConvertInstance converts an instantiated cell to a QMASM macro instantiation.
// This function returns a slice rather than an individual QmasmMacroUse because
// it may need to return an empty slice (in the case of VCC and GND).
func ConvertInstance(inst EdifList, i2n map[EdifSymbol]EdifString) []QmasmCode {
	// Extract the instantiation name.
	code := make([]QmasmCode, 0, 1)
	instSym, comment := nameAndComment(inst[1])
	if instSym == "GND" || instSym == "VCC" {
		return nil // GND and VCC are treated specially.
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

	// Keep track of whether the macro represents a flip-flop.
	isFlipFlop[instSym] = false
	if macroName == "DFF_P" || macroName == "DFF_N" {
		isFlipFlop[instSym] = true
	}

	// Construct and return a macro instantiation.
	code = append(code, QmasmMacroUse{
		MacroName: string(macroName),
		UseNames:  []string{"$" + string(instSym)},
		Comment:   comment,
	})
	return code
}

// PortRefToString converts an EDIF portRef to a string.  This is a helper
// function for ConvertNet.
func PortRefToString(pRef EdifList) string {
	// We can handle only 2- or 3-element portRefs.
	nParts := len(pRef)
	if nParts != 2 && nParts != 3 {
		notify.Fatalf("Expected 2 or 3 elements in a portRef; saw %v", pRef)
	}

	// The first element after "portRef" is the port name.
	var pName string
	switch pRef[1].Type() {
	case Symbol:
		// Single-bit
		pName = string(AsSymbol(pRef[1]))

	case List:
		// Index into a multi-bit port.  Return as "symbol[port]".
		memb := AsList(pRef[1], 3, "member")
		base := AsSymbol(memb[1])
		idx := AsInteger(memb[2])
		pName = fmt.Sprintf("%s[%d]", base, idx)

	default:
		notify.Fatalf("Expected a symbol or list in portRef but saw %v", pRef)
	}

	// If provided, the second element after "portRef" is the cell the port
	// belongs to.
	if nParts <= 2 {
		return pName
	}
	instRef := AsList(pRef[2], 2, "instanceRef")
	return "$" + string(AsSymbol(instRef[1])) + "." + pName
}

// PortRefFlipFlopPort takes an EDIF portRef and returns its port name if the
// cell represents a flip-flip or the empty string otherwise.  This is a helper
// function for ConvertNet.
func PortRefFlipFlopPort(pRef EdifList) string {
	// We can handle only 2- or 3-element portRefs.
	nParts := len(pRef)
	if nParts != 2 && nParts != 3 {
		notify.Fatalf("Expected 2 or 3 elements in a portRef; saw %v", pRef)
	}

	// The first element after "portRef" is the port name.
	var pName string
	switch pRef[1].Type() {
	case Symbol:
		// Single-bit
		pName = string(AsSymbol(pRef[1]))

	case List:
		// Index into a multi-bit port.  Flip-flops never use this
		// feature.
		return ""

	default:
		notify.Fatalf("Expected a symbol or list in portRef but saw %v", pRef)
	}

	// If provided, the second element after "portRef" is the cell the port
	// belongs to.
	if nParts <= 2 {
		return ""
	}
	instRef := AsList(pRef[2], 2, "instanceRef")
	instSym := AsSymbol(instRef[1])
	if isFlipFlop[instSym] {
		return pName
	}
	return ""
}

// arrayIndexRe matches an array element such as "foo[123]".  The first
// capturing group is the array name, and the second is the index value.
var arrayIndexRe = regexp.MustCompile(`^([^\[\]]+)\[(\d+)\]$`)

// needsRenaming determines heuristically if a symbol and a comment refer to
// different bits within the same bit vector.  This indicates an endianness
// mismatch between the HDL and the netlist representation.  We will water
// rename the symbol to match the HDL version.
func needsRenaming(s, c string) bool {
	if c == s {
		return false
	}
	cai := arrayIndexRe.FindStringSubmatch(c)
	if cai == nil {
		return false
	}
	sai := arrayIndexRe.FindStringSubmatch(s)
	if sai == nil {
		return false
	}
	return cai[1] == sai[1] && cai[2] != sai[2] // Same name, different indices
}

// ConvertNet converts an EDIF net to a QMASM chain ("=").
func ConvertNet(net EdifList, iface map[EdifSymbol]Empty) []QmasmCode {
	// Keep track of port names and flip-flop status.
	type PortInfo struct {
		Name   string
		FFPort string
	}
	pInfo := make([]PortInfo, 0, 2)

	// Determine the name of each port.
	for _, pRef := range net.NestedSublistsByName([]EdifSymbol{
		"joined",
		"portRef",
	}) {
		pInfo = append(pInfo, PortInfo{
			Name:   PortRefToString(pRef),
			FFPort: PortRefFlipFlopPort(pRef),
		})
	}
	if len(pInfo) < 2 {
		// I don't know what a single-portRef net is supposed to do so
		// I'm guessing we can ignore it.
		return nil
	}

	// Treat a renamed net as a comment.
	_, comment := nameAndComment(net[1])

	// Rename EDIF array accesses to HDL array accesses to account for
	// endianness differences.
	nPorts := len(pInfo)
	code := make([]QmasmCode, 0, (nPorts*(nPorts-1))/2)
	special := map[string]bool{
		"$GND.G": false,
		"$VCC.P": true,
	}
	for i := 0; i < nPorts; i++ {
		_, iPinned := special[pInfo[i].Name]
		iPrefix := ""
		if pInfo[i].FFPort == "Q" {
			iPrefix = "!next."
		}
		if iPinned {
			continue
		}
		iName := iPrefix + pInfo[i].Name
		if needsRenaming(iName, comment) {
			code = append(code, QmasmRename{
				Before: []string{iName},
				After:  []string{comment},
			})
		}
	}
	if len(code) > 0 {
		// At least one symbol was renamed.  Alter the comment to make
		// it more user-friendly.
		comment = comment + " in the HDL"
	}

	// Return one or more QMASM chains/pins.
	for i := 0; i < nPorts-1; i++ {
		for j := i + 1; j < nPorts; j++ {
			iVal, iPinned := special[pInfo[i].Name]
			iPrefix := ""
			if pInfo[i].FFPort == "Q" {
				iPrefix = "!next."
			}
			jVal, jPinned := special[pInfo[j].Name]
			jPrefix := ""
			if pInfo[j].FFPort == "Q" {
				jPrefix = "!next."
			}
			switch {
			case !iPinned && !jPinned:
				// Neither port is VCC or GND.
				iName := iPrefix + pInfo[i].Name
				jName := jPrefix + pInfo[j].Name
				if iName == jName {
					// I'm not convinced we'll ever get
					// here in practice.
					continue
				}
				code = append(code, QmasmChain{
					Var:     [2]string{iName, jName},
					Comment: comment,
				})

			case iPinned && !jPinned:
				// Only port i is VCC or GND.
				code = append(code, QmasmPin{
					Var:     pInfo[i].Name,
					Value:   iVal,
					Comment: comment,
				})

			case !iPinned && jPinned:
				// Only port j is VCC or GND.
				code = append(code, QmasmPin{
					Var:     pInfo[j].Name,
					Value:   jVal,
					Comment: comment,
				})

			default:
				notify.Fatalf("Unexpected connection in net %v", net)
			}
		}
	}
	return code
}

// ParseInterface extracts a cell interface and parses it into a set of port
// names.
func ParseInterface(cell EdifList) map[EdifSymbol]Empty {
	// Find the interface.
	ifs := cell.NestedSublistsByName([]EdifSymbol{"view", "interface"})
	if len(ifs) != 1 {
		notify.Fatalf("Expected exactly one interface; saw %d", len(ifs))
	}

	// Process each port in the interface in turn.
	pNames := make(map[EdifSymbol]Empty, len(ifs[0])-1)
	for _, p := range ifs[0][1:] {
		port := AsList(p, 3, "port")
		switch port[1].Type() {
		case Symbol:
			// Single bit
			pNames[AsSymbol(port[1])] = Empty{}

		case List:
			pList := port[1].(EdifList)
			switch AsSymbol(pList[0]) {
			case "array":
				// Array of bits, zero-based.
				array := AsList(port[1], 3, "array")
				aLen := int(AsInteger(array[2]))
				bSym, base := nameAndComment(array[1])
				if base == "" {
					base = string(bSym)
				}
				for i := 0; i < aLen; i++ {
					sym := fmt.Sprintf("%s[%d]", base, i)
					pNames[EdifSymbol(sym)] = Empty{}
				}

			case "rename":
				// Renamed single bit
				sym, comment := nameAndComment(port[1])
				if comment != "" {
					sym = EdifSymbol(sym)
				}
				pNames[sym] = Empty{}

			default:
				notify.Fatalf("Failed to parse a port list of type %q", AsSymbol(pList[0]))
			}

		default:
			notify.Fatalf("Expected a symbol or list in port but saw %v", port)
		}
	}
	return pNames
}

// ConvertCell converts a user-defined cell to a QMASM macro definition.
func ConvertCell(cell EdifList, i2n map[EdifSymbol]EdifString) QmasmMacroDef {
	// Ensure the cell looks at least a little like what we expect.
	if len(cell) < 3 {
		notify.Fatalf("Cell %v contains too few components", cell)
	}

	// Extract the cell's external interface.
	iface := ParseInterface(cell)

	// Instantiate all the other cells used by the current cell.
	code := make(QmasmCodeList, 0, 32)
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
		code = append(code, ConvertNet(net, iface)...)
	}

	// Sort the code, wrap it in a QMASM macro definition, and return it.
	code = code.SortAndMerge()
	cName, cComment := nameAndComment(cell[1])
	return QmasmMacroDef{
		Name:    string(cName),
		Body:    code,
		Comment: cComment,
	}
}

// ConvertDesign converts a design to a QMASM macro instantiation.
func ConvertDesign(des EdifList, nCycles uint, noTop bool) QmasmMacroUse {
	if len(des) != 3 {
		notify.Fatalf("Expected a design to contain exactly 3 elements but saw %v", des)
	}
	cRef := AsList(des[2], 3, "cellRef")
	name, comment := nameAndComment(des[1])
	if comment != "" {
		name = EdifSymbol(comment)
	}
	uNames := make([]string, nCycles)
	switch {
	case nCycles > 1:
		for i := range uNames {
			uNames[i] = fmt.Sprintf("%s@%d", name, i)
		}
	case noTop:
		uNames = make([]string, 0)
	default:
		uNames[0] = string(name)
	}
	return QmasmMacroUse{
		MacroName: string(AsSymbol(cRef[1])),
		UseNames:  uNames,
	}
}

// ConvertLibrary converts a user-defined cell library to QMASM macro
// definitions.
func ConvertLibrary(lib EdifList, i2n map[EdifSymbol]EdifString) []QmasmCode {
	// Iterate over each cell.
	code := make([]QmasmCode, 0, 32)
	for _, cell := range lib.SublistsByName("cell") {
		code = append(code, QmasmBlank{})
		code = append(code, ConvertCell(cell, i2n))
	}
	return code
}

// ConvertEdifToQmasm takes an EDIF s-expression and returns a list of QMASM
// statements.
func ConvertEdifToQmasm(s EdifSExp, nCycles uint, noTop bool) []QmasmCode {
	// Produce a QMASM header block.
	code := make([]QmasmCode, 0, 128)
	code = append(code, ConvertMetadata(s)...)
	code = append(code, QmasmBlank{})
	code = append(code, QmasmInclude{File: "stdcell"})

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
		code = append(code, ConvertLibrary(lib, idToName)...)
	}

	// Convert each design in turn.
	for _, des := range slst.SublistsByName("design") {
		code = append(code, QmasmBlank{})
		code = append(code, ConvertDesign(des, nCycles, noTop))
	}

	return code
}
