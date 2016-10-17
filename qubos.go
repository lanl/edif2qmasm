// This file is part of edif2qubo.  It provides functions for manipulating
// QUBOs.

package main

import (
	"fmt"
	"io"
)

// A QuboEltType describes the type of an element of a QUBO problem.
type QuboEltType int

// These are the types of QUBO elements we support.
const (
	Point   QuboEltType = iota // Point weight
	Coupler                    // Coupler strength
	Chain                      // Chain between two symbols
	Alias                      // Equivalence of two symbols
)

// A QuboElt is one element of a QUBO problem.  A list of QuboElts defines a
// complete QUBO problem.
type QuboElt struct {
	Type   QuboEltType // How to interpret the following fields
	X      EdifSymbol  // First (or only symbol)
	Y      EdifSymbol  // Second symbol (if needed)
	Weight int         // Point weight or coupler strength
}

// A Qubo is an ordered list of (symbolic) weights and couplers.
type Qubo []QuboElt

// OutputText outputs a QUBO with one weight or coupler strength per line.
func (q Qubo) OutputText(w io.Writer) error {
	var err error
	for _, elt := range q {
		switch elt.Type {
		case Point:
			_, err = fmt.Fprintf(w, "%s %d\n", elt.X, elt.Weight)
		case Coupler:
			_, err = fmt.Fprintf(w, "%s %s %d\n", elt.X, elt.Y, elt.Weight)
		case Chain:
			_, err = fmt.Fprintf(w, "%s = %s\n", elt.X, elt.Y)
		case Alias:
			_, err = fmt.Fprintf(w, "%s == %s\n", elt.X, elt.Y)
		default:
			notify.Fatalf("Unexpected element type %d", elt.Type)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
