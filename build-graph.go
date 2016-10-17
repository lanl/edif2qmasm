// This file is part of edif2qasm.  It constructs an internal graph
// representation of a netlist.

package main

import "fmt"

// CellAndCode is a tuple that stores a Cell and the EDIF code that produced it.
type CellAndCode struct {
	C Cell
	E EdifList
}

// A CellLibrary maps from a cell name to a Cell and the associated EDIF code.
type CellLibrary map[EdifSymbol]CellAndCode

// asSymbol asserts that an s-expression is a symbol, returning it if it is but
// aborting the progam if it's not.
func asSymbol(s EdifSExp) EdifSymbol {
	if s.Type() != Symbol {
		notify.Fatalf("Expected a Symbol but received a %s (%v)", s.Type(), s)
	}
	return s.(EdifSymbol)
}

// asString asserts that an s-expression is a string, aborting the progam if
// not.
func asString(s EdifSExp) EdifString {
	if s.Type() != String {
		notify.Fatalf("Expected a String but received a %s (%v)", s.Type(), s)
	}
	return s.(EdifString)
}

// asInteger asserts that an s-expression is an integer, aborting the progam if
// not.
func asInteger(s EdifSExp) EdifInteger {
	if s.Type() != Integer {
		notify.Fatalf("Expected an Integer but received a %s (%v)", s.Type(), s)
	}
	return s.(EdifInteger)
}

// asList asserts that an s-expression is a list, contains a minimum number of
// elements, and begins with a given keyword.  If so, it returns the list.
// Otherwise, it aborts the progam.
func asList(s EdifSExp, minElts int, keyw EdifSymbol) EdifList {
	if s.Type() != List {
		notify.Fatalf("Expected a [%s ...] List but received a %s (%v)", keyw, s.Type(), s)
	}
	lst := s.(EdifList)
	if len(lst) < minElts {
		notify.Fatalf("Expected list %v to contain at least %d elements", lst, minElts)
	}
	if minElts > 0 && keyw != "" {
		k := asSymbol(lst[0])
		if k != keyw {
			notify.Fatalf("Expected keyword %q but got %q", keyw, k)
		}
	}
	return lst
}

// SublistsByName returns all direct sublists of the form "(<keyw> ...)".
func SublistsByName(list EdifList, keyw EdifSymbol) []EdifList {
	eltList := make([]EdifList, 0)
	for _, elt := range list {
		if elt.Type() != List {
			continue
		}
		subList := asList(elt, 0, "")
		if len(subList) == 0 || subList[0].Type() != Symbol {
			continue
		}
		if asSymbol(subList[0]) != keyw {
			continue
		}
		eltList = append(eltList, subList)
	}
	return eltList
}

// NestedSublistsByName returns all sublists of the form "(<keyw> ...)"  that
// are nested within "(<keyw1> (<keyw2> (<keyw3> ...)))".
func NestedSublistsByName(list EdifList, kws []EdifSymbol) []EdifList {
	switch len(kws) {
	case 0:
		return nil
	case 1:
		return SublistsByName(list, kws[0])
	default:
		subLists := make([]EdifList, 0)
		for _, outerList := range SublistsByName(list, kws[0]) {
			innerLists := NestedSublistsByName(outerList, kws[1:])
			subLists = append(subLists, innerLists...)
		}
		return subLists
	}
}

// FirstSymbol returns the first raw symbol within a list.
func FirstSymbol(list EdifList) EdifSymbol {
	for _, elt := range list {
		if elt.Type() == Symbol {
			return asSymbol(elt)
		}
	}
	return ""
}

// SymbolAndAlias parses either list of the form "(rename <symbol> <string>)" or
// a single symbol into a {<symbol>, <string>} pair.  It aborts on any parse
// error.
func SymbolAndAlias(elt EdifSExp) (EdifSymbol, EdifString) {
	if elt.Type() == Symbol {
		sym := asSymbol(elt)
		return sym, EdifString(sym)
	}
	lst := asList(elt, 3, "rename")
	return asSymbol(lst[1]), asString(lst[2])
}

// ProcessCell creates a Cell based on a "(cell ...)" s-expression and returns
// its symbolic ID and textual name.
func ProcessCell(userDef bool, cList EdifList) (Cell, EdifSymbol, EdifString) {
	// Extract the cell name(s).
	if len(cList) < 3 {
		notify.Fatal("Cell contains too few elements")
	}
	cName, cAlias := SymbolAndAlias(cList[1])

	// Create a cell of the appropriate concrete type.
	var cell Cell
	if userDef {
		cell = &UserCell{}
	} else {
		switch cAlias {
		case "GND":
			cell = &GndCell{}
		case "VCC":
			cell = &VccCell{}
		case "$_NOT_":
			cell = &NotCell{}
		case "$_AND_":
			cell = &AndCell{}
		case "$_MUX_":
			cell = &MuxCell{}
		default:
			notify.Fatalf("No cell defined for name(s) %v", cAlias)
		}
	}

	// Extract the ports and assign them to the cell.
	for _, port := range NestedSublistsByName(cList,
		[]EdifSymbol{"view", "interface", "port"}) {
		// We expect the port to have the form "(port <name>
		// (direction <INPUT|OUTPUT>)".
		if len(port) != 3 {
			notify.Fatalf("Unable to parse a %d-element port", len(port))
		}

		// Extract the port name(s)
		pNames := make([]EdifSymbol, 0, 1)
		if port[1].Type() == Symbol {
			pNames = append(pNames, asSymbol(port[1]))
		} else {
			array := asList(port[1], 3, "array")
			base := asSymbol(array[1])
			width := asInteger(array[2])
			for i := EdifInteger(0); i < width; i++ {
				name := fmt.Sprintf("%s[%d]", base, i)
				pNames = append(pNames, EdifSymbol(name))
			}
		}

		// Extract the port direction.
		dirList := asList(port[2], 2, "direction")
		switch asSymbol(dirList[1]) {
		case "INPUT":
			for _, p := range pNames {
				cell.CreateInput(p)
			}
		case "OUTPUT":
			for _, p := range pNames {
				cell.CreateOutput(p)
			}
		default:
			notify.Fatalf("Expected INPUT or OUTPUT but got %q", asSymbol(dirList[1]))
		}
	}

	// Return the constructed cell and its names.
	return cell, cName, cAlias
}

// ProcessCells defines logic cells according to a "(external|library ...)"
// s-expression.
func ProcessCells(external bool, libList EdifList) (EdifSymbol, CellLibrary) {
	libName := FirstSymbol(libList[1:])
	lib := make(CellLibrary)
	for _, cList := range SublistsByName(libList[1:], "cell") {
		cell, cName, _ := ProcessCell(external, cList)
		lib[cName] = CellAndCode{C: cell, E: cList}
	}
	return libName, lib
}

// FullyInstantiate instantiates a cell and, recursively, every cell that it
// instantiates.
func FullyInstantiate(prefix EdifString, cellInfo CellAndCode,
	defLibName EdifSymbol, cellLibs map[EdifSymbol]CellLibrary) Instance {
	// Extract the cell's "(contents ...)" s-expression.
	cLists := NestedSublistsByName(cellInfo.E, []EdifSymbol{"view", "contents"})
	if len(cLists) > 1 {
		notify.Fatalf("Expected only one set of contents per cell; received %d", len(cLists))
	}
	newInstance := cellInfo.C.MakeInstance() // New instance to return
	if len(cLists) == 0 {
		// No children --> no more work to do
		return newInstance
	}
	contList := asList(cLists[0], 0, "contents")

	// Recursively instantiate each cell referred to in the contents.  Add
	// each child instance to our instance.
	children := make(map[EdifSymbol]Instance) // Cache of our child instances
	for _, inst := range SublistsByName(contList, "instance") {
		instList := asList(inst, 3, "instance")
		iName, _ := SymbolAndAlias(instList[1])
		vrefList := asList(instList[2], 3, "viewRef")
		crefList := asList(vrefList[2], 2, "cellRef")
		lName := defLibName // Name of library that defines the cell to instantiate
		if len(crefList) >= 3 {
			// Reference to a cell in a difference library.
			libRef := asList(crefList[2], 2, "libraryRef")
			lName = asSymbol(libRef[1])
			if _, ok := cellLibs[lName]; !ok {
				notify.Fatalf("Failed to find library %s", lName)
			}
		}
		cName := asSymbol(crefList[1])
		cInfo, ok := cellLibs[lName][cName]
		if !ok {
			notify.Fatalf("Failed to find cell %s", cName)
		}
		newPrefix := prefix + "." + EdifString(iName)
		child := FullyInstantiate(newPrefix, cInfo, lName, cellLibs)
		child.SetName(prefix + "." + EdifString(iName))
		newInstance.IncludeChild(iName, child)
		children[iName] = child
	}

	// Define a helper function that converts portRefs to {instance, port}
	// pairs.
	parsePortRef := func(pRef EdifSExp) (Instance, EdifSymbol) {
		// Find the port name, which can be either a symbol or an array
		// element.
		pr := asList(pRef, 2, "portRef")
		var pName EdifSymbol // Port name
		switch pr[1].Type() {
		case Symbol:
			pName = asSymbol(pr[1])
		case List:
			mem := asList(pr[1], 3, "member")
			v := asSymbol(mem[1])
			idx := asInteger(mem[2])
			pName = EdifSymbol(fmt.Sprintf("%s[%d]", v, idx))
		default:
			notify.Fatalf("Failed to parse port %v", pr)
		}

		// Find the cell instance that to which the port belongs.
		var pInst Instance // Instance containing pName
		if len(pr) <= 2 {
			// Instance we just created
			pInst = newInstance
		} else {
			// Child instance
			instRef := asList(pr[2], 2, "instanceRef")
			pInst = children[asSymbol(instRef[1])]
		}

		// Return what we found.
		return pInst, pName
	}

	// Wire up our child instances.
	for _, net := range SublistsByName(contList, "net") {
		if len(net) != 3 {
			notify.Fatalf("Net %v contains %d elements, not 3", net, len(net))
		}
		joined := asList(net[2], 3, "joined")
		inst1, port1 := parsePortRef(joined[len(joined)-1])
		type1, ok := inst1.TypeOfPort(port1)
		if !ok {
			notify.Fatalf("Failed to find port %s", port1)
		}
		for _, otherRef := range joined[1 : len(joined)-1] {
			// Connect two ports together.
			inst2, port2 := parsePortRef(otherRef)
			type2, ok := inst2.TypeOfPort(port2)
			if !ok {
				notify.Fatalf("Failed to find port %s", port2)
			}
			if type1 == type2 {
				// Represent input-input and output-output with
				// aliases, not connections.  We alias from
				// parent to child.
				if inst1 == newInstance {
					inst1.AliasPort(port1, inst2, port2)
				} else {
					inst2.AliasPort(port2, inst1, port1)
				}
			} else {
				// Represent output-input with a connection.
				// Input-output fails.
				inst1.ConnectPort(port1, inst2, port2)
			}
		}
	}
	return newInstance
}

// ConvertEdifToNetlist converts a parsed EDIF file to an internal graph
// representation.
func ConvertEdifToNetlist(edif EdifSExp) map[EdifSymbol]Instance {
	// Process the top level of the s-expression, which should be of the
	// form "(edif ...)" and contain, at a minimum a symbolic name for the
	// design, an "(external ...)" list of cells, a "(library ...)"  list
	// of cells and instantiations, and a "(design ...)" that specifies the
	// top-level design.
	eList := asList(edif, 5, "edif")

	// Process all externals to construct a list of externally defined
	// cells.
	allCells := make(map[EdifSymbol]CellLibrary) // Map from library name to cell name to cell information
	for _, extList := range SublistsByName(eList, "external") {
		libName, lib := ProcessCells(false, extList)
		allCells[libName] = lib
	}

	// Process all libraries to construct a list of user-defined cells.
	for _, libList := range SublistsByName(eList, "library") {
		libName, lib := ProcessCells(true, libList)
		allCells[libName] = lib
	}

	// Instantiate the top-level design(s).
	designs := make(map[EdifSymbol]Instance)
	for _, desList := range SublistsByName(eList, "design") {
		if len(desList) != 3 {
			notify.Fatalf("Design list %v has %d element(s), not 3", desList, len(desList))
		}
		dName := asSymbol(desList[1])
		cRef := asList(desList[2], 3, "cellRef")
		cName := asSymbol(cRef[1])
		lRef := asList(cRef[2], 2, "libraryRef")
		lName := asSymbol(lRef[1])
		lib, ok := allCells[lName]
		if !ok {
			notify.Fatalf("Failed to find library %s", lName)
		}
		cInfo, ok := lib[cName]
		if !ok {
			notify.Fatalf("Failed to find cell %s", cName)
		}
		inst := FullyInstantiate(EdifString(dName), cInfo, lName, allCells)
		inst.SetName(EdifString(dName))
		designs[dName] = inst
	}
	return designs
}
