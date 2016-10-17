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

// ConvertEdifToQasm takes an EDIF s-expression and returns a list of QASM
// statements.
func ConvertEdifToQasm(s EdifSExp) []QasmCode {
	code := make([]QasmCode, 0, 128)
	code = append(code, ConvertMetadata(s)...)
	return code
}
