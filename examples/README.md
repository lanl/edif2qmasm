edif2qmasm Examples
===================

This directory contains a few examples written in [Verilog](https://en.wikipedia.org/wiki/Verilog) and a `Makefile` that compiles them to [EDIF](https://en.wikipedia.org/wiki/EDIF) netlists using the [Yosys Open SYnthesis Suite](http://www.clifford.at/yosys/).  Start by running
```bash
make
```
to produce a `.edif` file from each `.v` file or use your favorite hardware-synthesis tool to perform the equivalent operation.

`edif2qmasm` currently supports only a handful of gates (defined in [stdcell.qmasm](https://github.com/losalamos/edif2qmasm/blob/master/stdcell.qmasm)) so all designs must be compiled to use only those gates.

* 1-input: NOT
* 2-input: AND, OR, XOR
* 3-input: MUX

The rest of this document describes each of the examples in turn.

circsat
-------

[`circsat.v`](https://github.com/losalamos/edif2qmasm/blob/master/examples/circsat.v) implements an arbitrary circuit, specifically the one presented in [Cormen, Leiserson, Rivest, and Stein's Algorithms textbook](https://mitpress.mit.edu/books/introduction-algorithms) in its discussion of circuit-satisfiability problems.  The Verilog module's inputs are named *a*, *b*, and *c*, and the sole output is named *y*.  All are single bits.  Internally, inputs, outputs, and intermediate values are named *x*[1]…*x*[10].

The goal of this example is to run the circuit *backward* to find out what set of inputs produces an output of *true*.  This is a classic NP-complete problem.  Here's how to run it on a D-Wave system using `edif2qmasm` and `qmasm`:
```bash
$ edif2qmasm circsat.edif | qmasm --run --pin="circsat.y := true"
# circsat.a --> 1033
# circsat.b --> 1035
# circsat.c --> 936 941
# circsat.y --> 733
Solution #1 (energy = -50.25, tally = 3):

    Name(s)       Spin  Boolean
    ------------  ----  --------
    circsat.a       +1  True
    circsat.b       +1  True
    circsat.c       -1  False
    circsat.y       +1  True
```

Note that we pinned the output to *true* while leaving the other parameters unspecified.  Doing so led the D-Wave to search for a minimum-energy solution subject to that constraint, and it found only {*true*, *true*, *false*}.  The circuit specified by `circsat.v` is small enough that we can fully evaluate all possibilities and see that this in indeed the sole solution:

| *a* = *x*[1] | *b* = *x*[2] | *c* = *x*[3]  | *x*[4]  | *x*[5]  | *x*[6]  | *x*[7]  | *x*[8]  | *x*[9]  | *y* = *x*[10] |
| ------------ | ------------ | ------------- | ------- | ------- | ------- | ------- | ------- | ------- | ------------- |
| FALSE        | FALSE        | FALSE         | TRUE    | FALSE   | FALSE   | FALSE   | FALSE   | FALSE   | FALSE         |
| FALSE        | FALSE        | TRUE          | FALSE   | FALSE   | TRUE    | FALSE   | TRUE    | TRUE    | FALSE         |
| FALSE        | TRUE         | FALSE         | TRUE    | TRUE    | FALSE   | FALSE   | TRUE    | FALSE   | FALSE         |
| FALSE        | TRUE         | TRUE          | FALSE   | TRUE    | TRUE    | FALSE   | TRUE    | TRUE    | FALSE         |
| TRUE         | FALSE        | FALSE         | TRUE    | TRUE    | FALSE   | FALSE   | TRUE    | FALSE   | FALSE         |
| TRUE         | FALSE        | TRUE          | FALSE   | TRUE    | TRUE    | FALSE   | TRUE    | TRUE    | FALSE         |
| TRUE         | TRUE         | FALSE         | TRUE    | TRUE    | FALSE   | TRUE    | TRUE    | TRUE    | TRUE          |
| TRUE         | TRUE         | TRUE          | FALSE   | TRUE    | TRUE    | FALSE   | TRUE    | TRUE    | FALSE         |
