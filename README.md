edif2qmasm
==========

[![Build Status](https://travis-ci.org/lanl/edif2qmasm.svg?branch=master)](https://travis-ci.org/lanl/edif2qmasm) [![Go Report Card](https://goreportcard.com/badge/github.com/lanl/edif2qmasm)](https://goreportcard.com/report/github.com/lanl/edif2qmasm)

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

More precisely, `edif2qmasm` converts from the [EDIF](https://en.wikipedia.org/wiki/EDIF) netlist format, which can be output by various synthesis tools, to the [QMASM](https://github.com/lanl/qmasm) quantum macro assembly language. To date, `edif2qmasm` has been tested only with Verilog because there exist open-source compilers that convert Verilog to EDIF, and I don't know of an equivalent open-source tool that can convert VHDL to EDIF.

Documentation
-------------

In addition to this file, two other sources of documentation are

* [the `edif2qmasm` manual page](https://github.com/lanl/edif2qmasm/blob/master/edif2qmasm.rst) and
* [the `README.md` file in the `examples` directory](https://github.com/lanl/edif2qmasm/blob/master/examples/README.md).

There also exists a peer-reviewed academic publication on `edif2qmasm` that describes the entire process of compiling Verilog programs to a D-Wave Hamiltonian function, explains how `edif2qmasm`'s standard-cell library was constructed, presents some use cases, and even includes a bit of analysis:

> Scott Pakin.  "Targeting Classical Code to a Quantum Annealer".  In <em>Proceedings of the Twenty-Fourth International Conference on Architectural Support for Programming Languages and Operating Systems (ASPLOS 2019)</em>, 13–17 April 2019, Providence, Rhode Island, USA, pp. 529–543.  ACM, New York, New York, USA.  ISBN: 978-1-4503-6240-5, DOI: [10.1145/3297858.3304071](https://doi.org/10.1145/3297858.3304071).

Associated with the above is a [2-minute "lightning talk" video](https://youtu.be/jtFsujUM-4Q), which is essentially an advertisement to attend the talk (which was 16 April 2019) and read the paper.

Installation
------------

`edif2qmasm` is written in [Go](https://go.dev/) and therefore depends upon a Go compiler to build.  See [the INSTALL file](https://github.com/lanl/edif2qmasm/blob/master/INSTALL.md) for build instructions.

Usage
-----

`edif2qmasm` usage is straightforward:
```bash
edif2qmasm -o myfile.qmasm myfile.edif
```
If no input file is specified, `edif2qmasm` will read from the standard input device.  Run `edif2qmasm --help` for a list of available command-line options.

To run the generated code with QMASM, you'll need to point it to the `edif2qmasm` standard-cell library.  In Bash, enter
```bash
export QMASMPATH=/usr/local/share/edif2qmasm:$QMASMPATH
```
replacing `/usr/local` with whatever installation prefix you used.

Limitations
-----------

`edif2qmasm` has only limited support for sequential logic.  Sequential logic is implemented by replicating the entire circuit once per clock cycle for a compile-time specified number of clock cycles (cf. the `--cycles` command-line option).  Clocked flip-flops are supported, but unclocked latches are not.

The resulting QMASM programs are not very robust in that the minimum-energy solutions do not consistently represent a correct execution when run on D-Wave hardware.  Running QMASM with `--postproc=opt` helps substantially.  Other suggestions on how to improve robustness are welcome.

License
-------

`edif2qmasm` is provided under a BSD-ish license with a "modifications must be indicated" clause.  See [the LICENSE file](https://github.com/lanl/edif2qmasm/blob/master/LICENSE.md) for the full text.

`edif2qmasm` is part of the Hybrid Quantum-Classical Computing suite, known internally as LA-CC-16-032.

Author
------

Scott Pakin, <pakin@lanl.gov>
