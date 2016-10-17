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
func ProcessExternalLib(s EdifSExp) map[EdifSymbol]EdifString {
	idToName := make(map[EdifSymbol]EdifString, 8)
	ext := AsList(s, 2, "external")
	for _, cell := range ext.SublistsByName("cell") {
		cnm := cell[1]
		if cnm.Type() != List {
			continue // Symbols don't need to be mapped.
		}
		rnm := AsList(cnm, 2, "rename")
		idToName[AsSymbol(rnm[1])] = CanonicalizeCellName(AsString(rnm[2]))
	}
	return idToName
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

	return code
}
