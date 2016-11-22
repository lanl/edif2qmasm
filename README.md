edif2qmasm
==========

Description
-----------

`edif2qmasm` makes it possible to run [Verilog](https://en.wikipedia.org/wiki/Verilog) or [VHDL](https://en.wikipedia.org/wiki/VHDL) programs on a [D-Wave quantum annealer](http://www.dwavesys.com/).

*Why quantum annealing?*  The primary reason to target a quantum annealer—really, any hardware annealer—is that programs can be run in either the forward or backward direction.  One can in fact specify any combination of inputs and outputs and solve for the other, unspecified, values.  This benefits both

* expressiveness, as some programs are easier to write in the _A_ → _B_ direction than in the _B_ → _A_ direction, and
* performance, as some classical algorithms run quickly in the _A_ → _B_ direction but slowly in the _B_ → _A_ direction.  Consider verifying a solution to an [NP-complete](https://en.wikipedia.org/wiki/NP-completeness) problem (fast) versus producing a solution to an NP-complete problem (slow).

*Why Verilog/VHDL?*  Some of the advantages of using a hardware-description language as a D-Wave programming language are that it

* supports basic programming-language features such as conditionals, loops, multi-bit constants and variables, assignments, arithmetic operations, and modules,
* provides precise control over bit widths, which reduces the number of wasted qubits (a precious resource in contemporary D-Wave systems),
* enables exploiting the code optimizations and debugging support provided by synthesis tools.

More precisely, `edif2qmasm` converts from the [EDIF](https://en.wikipedia.org/wiki/EDIF) netlist format, which can be output by various synthesis tools, to the [QMASM](https://github.com/losalamos/qmasm) quantum macro assembly language. To date, `edif2qmasm` has been tested only with Verilog because there exist open-source compilers that convert Verilog to EDIF, and I don't know of an equivalent open-source tool that can convert VHDL to EDIF.

Installation
------------

`edif2qmasm` is written in [Go](https://golang.org/) and therefore depends upon a Go compiler to build.  The program's build process further requires [`goimports`](https://godoc.org/golang.org/x/tools/cmd/goimports), [`stringer`](https://godoc.org/golang.org/x/tools/cmd/stringer), and the [Pigeon parser generator](https://godoc.org/github.com/PuerkitoBio/pigeon).  These can be installed from the command line with
```bash
go get golang.org/x/tools/cmd/goimports
go get golang.org/x/tools/cmd/stringer
go get github.com/PuerkitoBio/pigeon
```
Once these dependencies are satisfied, a simple
```bash
make
```
should build the `edif2qmasm` executable and
```bash
make install
```
should install it.  Specify the `prefix` variable to change the installation prefix from `/usr/local` to something else, e.g.,
```bash
make install prefix=$HOME
```

`edif2qmasm` is of limited use withut a compiler than can produce EDIF netlists and [QMASM](https://github.com/losalamos/qmasm), which executes the generated code on a D-Wave system.  To date, `edif2qmasm` has been tested only with the [Yosys Open SYnthesis Suite](http://www.clifford.at/yosys/), but reports of usage with other synthesis tools (successful or not) are welcome.  Note that QMASM relies on D-Wave's proprietary libraries to operate.

Usage
-----

`edif2qmasm` usage is straightforward:
```bash
edif2qmasm myfile.edif > myfile.qmasm
```
If no input file is specified, `edif2qmasm` will read from the standard input device.  There are not currently any command-line options.

To run the generated code with QMASM, you'll need to point it to the `edif2qmasm` standard-cell library.  In Bash, enter
```bash
export QMASMPATH=/usr/local/share/edif2qmasm:$QMASMPATH
```
replacing `/usr/local` with whatever installation prefix you used.

Limitations
-----------

Only combinational logic is supported by the current version of `edif2qmasm`.  I hope to add support for sequential logic in a future version.

The resulting QMASM programs are not very robust in that the minimum-energy solutions do not consistently represent a correct execution when run on D-Wave hardware.  Suggestions on how to improve robustness are welcome.

License
-------

`edif2qmasm` is provided under a BSD-ish license with a "modifications must be indicated" clause.  See [the LICENSE file](http://github.com/losalamos/edif2qmasm/blob/master/LICENSE.md) for the full text.

`edif2qmasm` is part of the Hybrid Quantum-Classical Computing suite, known internally as LA-CC-16-032.

Author
------

Scott Pakin, <pakin@lanl.gov>
