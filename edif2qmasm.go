// edif2qubits is a program that converts an EDIF file into LANL's
// QMASM format for execution on a quantum annealer.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

//go:generate bash -c "pigeon parse-edif.peg > parse-edif.tmp"
//go:generate bash -c "goimports parse-edif.tmp | gofmt > parse-edif.go"
//go:generate rm parse-edif.tmp
//go:generate stringer -type=SExpType edif.go edif2qmasm.go parse-edif.go qmasm.go  walk-sexp.go

var notify *log.Logger // Help notify the user of warnings and errors.

type Empty struct{} // Zero-byte type for defining and manipulating sets

func main() {
	// Parse the command line.
	var err error
	progName := path.Base(os.Args[0])
	notify = log.New(os.Stderr, progName+": ", 0)
	var r io.Reader
	switch len(os.Args) {
	case 1:
		r = os.Stdin
	case 2:
		f, err := os.Open(os.Args[1])
		if err != nil {
			notify.Fatal(err)
		}
		defer f.Close()
		r = f
	default:
		fmt.Fprintf(os.Stderr, "Usage: %s [<input.edif>]\n", progName)
		os.Exit(1)
	}

	// Parse the specified EDIF file into a top-level s-expression.
	parsed, err := ParseReader(progName, r)
	if err != nil {
		notify.Fatal(err)
	}
	top, ok := parsed.(EdifSExp)
	if !ok {
		notify.Fatalf("Failed to parse the input as an s-expression")
	}

	// Convert the s-expression to QMASM source code.
	code := ConvertEdifToQmasm(top)
	for _, q := range code {
		fmt.Printf("%s", q)
	}
}
