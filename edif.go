// This file is part of edif2qasm.  It abstracts an EDIF s-expression into a
// form that's easier to manipulate.

package main

// An SExpType represents the type of an EDIF s-expression.
type SExpType int

// These are the values an EDIF s-expression can take.
const (
	Symbol  SExpType = iota // Raw symbol
	String                  // Quoted string
	Integer                 // Integer
	List                    // List of any of the above or other lists
)

// An EdifSExp represents an EDIF s-expression.  It is effectively a union of
// various atom types and a list of atoms.
type EdifSExp interface {
	Type() SExpType // Return the type of this s-expression.
}

// An EdifSymbol represents a raw EDIF symbol.
type EdifSymbol string

// The type of an EdifSymbol is Symbol.
func (s EdifSymbol) Type() SExpType { return Symbol }

// An EdifString represents a quoted EDIF string.
type EdifString string

// The type of an EdifString is String.
func (s EdifString) Type() SExpType { return String }

// An EdifInteger represents an EDIF integer.
type EdifInteger int

// The type of an EdifInteger is Integer.
func (s EdifInteger) Type() SExpType { return Integer }

// An EdifList represents a list of EDIF s-expressions.
type EdifList []EdifSExp

// The type of an EdifList is List.
func (s EdifList) Type() SExpType { return List }

// AsSymbol asserts that an s-expression is a symbol, returning it if it is but
// aborting the program if it's not.
func AsSymbol(s EdifSExp) EdifSymbol {
	if s.Type() != Symbol {
		notify.Fatalf("Expected a Symbol but received a %s (%v)", s.Type(), s)
	}
	return s.(EdifSymbol)
}

// AsString asserts that an s-expression is a string, aborting the program if
// not.
func AsString(s EdifSExp) EdifString {
	if s.Type() != String {
		notify.Fatalf("Expected a String but received a %s (%v)", s.Type(), s)
	}
	return s.(EdifString)
}

// AsInteger asserts that an s-expression is an integer, aborting the program if
// not.
func AsInteger(s EdifSExp) EdifInteger {
	if s.Type() != Integer {
		notify.Fatalf("Expected an Integer but received a %s (%v)", s.Type(), s)
	}
	return s.(EdifInteger)
}

// AsList asserts that an s-expression is a list, contains a minimum number of
// elements, and begins with a given keyword.  If so, it returns the list.
// Otherwise, it aborts the program.
func AsList(s EdifSExp, minElts int, keyw EdifSymbol) EdifList {
	if s.Type() != List {
		notify.Fatalf("Expected a [%s ...] List but received a %s (%v)", keyw, s.Type(), s)
	}
	lst := s.(EdifList)
	if len(lst) < minElts {
		notify.Fatalf("Expected list %v to contain at least %d elements", lst, minElts)
	}
	if minElts > 0 && keyw != "" {
		k := AsSymbol(lst[0])
		if k != keyw {
			notify.Fatalf("Expected keyword %q but got %q", keyw, k)
		}
	}
	return lst
}

// SublistsByName returns all immediate sublists of the form "(<keyw> ...)".
func (l EdifList) SublistsByName(keyw EdifSymbol) []EdifList {
	eltList := make([]EdifList, 0)
	for _, elt := range l {
		if elt.Type() != List {
			continue
		}
		subList := AsList(elt, 0, "")
		if len(subList) == 0 || subList[0].Type() != Symbol {
			continue
		}
		if AsSymbol(subList[0]) != keyw {
			continue
		}
		eltList = append(eltList, subList)
	}
	return eltList
}

// NestedSublistsByName returns all sublists of the form "(<keyw> ...)"  that
// are nested within "(<keyw1> (<keyw2> (<keyw3> ...)))".
func (l EdifList) NestedSublistsByName(kws []EdifSymbol) []EdifList {
	switch len(kws) {
	case 0:
		return nil
	case 1:
		return l.SublistsByName(kws[0])
	default:
		subLists := make([]EdifList, 0)
		for _, outerList := range l.SublistsByName(kws[0]) {
			innerLists := outerList.NestedSublistsByName(kws[1:])
			subLists = append(subLists, innerLists...)
		}
		return subLists
	}
}
