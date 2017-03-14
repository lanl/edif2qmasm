/*
edif2qmasm is a program that converts an EDIF file into LANL's QMASM format
for execution on a quantum annealer.

What this means is that you can write a program in a hardware-description
language such as Verilog or VHDL then run it on a D-Wave quantum processing
unit (QPU).  The advantage of doing so is that the QPU does not distinguish
between inputs and outputs.  That is, programs can just as easily be run
"backward" (from outputs to inputs) as "forward" (from inputs to outputs)
or even a combination of the two.

For instance, a one-line program that assigns C = A*B can be given A and B
and produce their product, C; it can be given C and A and produce their
quotient, B; or it can be given the product C, and factor that into A and
B.  All of those variations consume the same amount of time but with the
caveat that a QPU is a stochastic device and is not guaranteed to produce
the same—or even a correct—answer every time.

Usage:

    edif2qmasm myfile.edif > myfile.qmasm

See https://github.com/lanl/edif2qmasm for more information.
*/
package main

import (
	"flag"
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

 // Empty is a zero-byte type for defining and manipulating sets.
type Empty struct{}

func main() {
	// Parse the command line.
	var err error
	progName := path.Base(os.Args[0])
	notify = log.New(os.Stderr, progName+": ", 0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [<option>]... [<input.edif>]\n", progName)
		flag.PrintDefaults()
	}
	nCycles := flag.Uint("cycles", 1, "number of clock cycles")
	outName := "-"
	flag.StringVar(&outName, "output", "-", "name of QMASM output file or \"-\" for stdout")
	flag.StringVar(&outName, "o", "-", "same as -output")
	flag.Parse()

	// Open the input file.
	var r io.Reader
	switch flag.NArg() {
	case 0:
		r = os.Stdin
	case 1:
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			notify.Fatal(err)
		}
		defer f.Close()
		r = f

	default:
		notify.Fatal("Too many input files")
	}

	// Open the output file.
	var out io.Writer
	if outName == "-" {
		out = os.Stdout
	} else {
		f, err := os.Create(outName)
		if err != nil {
			notify.Fatal(err)
		}
		defer f.Close()
		out = f
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
	code := ConvertEdifToQmasm(top, *nCycles)
	for _, q := range code {
		fmt.Fprintf(out, "%s", q)
	}
}
