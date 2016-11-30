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

mult
----

Excluding the comments, module definition, and parameter declarations, [`mult.v`](https://github.com/losalamos/edif2qmasm/blob/master/examples/mult.v) is a one-line Verilog program:
```Verilog
assign product = multiplicand * multiplier;
```

At the time of this writing, `mult.v` commonly produced incorrect solutions at the minimal-energy readings.  A debugging effort (with QMASM being the likely culprit) is underway.  Nevertheless, `mult.v` showcases the flexibility of `edif2qmasm`'s approach.  For starters, one can pin values to the inputs, say 3 and 5 (binary 011 and 101) to multiply them together:
```bash
$ edif2qmasm mult.edif | qmasm --run --pin="mult.multiplicand[2:0] := 011" --pin="mult.multiplier[2:0] := 101"
# mult.multiplicand[0] --> 233 239
# mult.multiplicand[1] --> 163 166
# mult.multiplicand[2] --> 388
# mult.multiplier[0] --> 192
# mult.multiplier[1] --> 716
# mult.multiplier[2] --> 158
# mult.product[0] --> 129
# mult.product[1] --> 839
# mult.product[2] --> 362
# mult.product[3] --> 25
# mult.product[4] --> 434
# mult.product[5] --> 404
Solution #1 (energy = -446.25, tally = 1):

    Name(s)                  Spin  Boolean
    -----------------------  ----  --------
    mult.multiplicand[0]       +1  True
    mult.multiplicand[1]       +1  True
    mult.multiplicand[2]       -1  False
    mult.multiplier[0]         +1  True
    mult.multiplier[1]         -1  False
    mult.multiplier[2]         +1  True
    mult.product[0]            +1  True
    mult.product[1]            +1  True
    mult.product[2]            +1  True
    mult.product[3]            -1  False
    mult.product[4]            -1  False
    mult.product[5]            -1  False
```

The `friendly-mult` script post-processes the above into a more human-readable form:
```bash
$ edif2qmasm mult.edif | qmasm --run --pin="mult.multiplicand[2:0] := 011" --pin="mult.multiplier[2:0] := 101" | ./friendly-mult
# mult.multiplicand[0] --> 998
# mult.multiplicand[1] --> 514 516
# mult.multiplicand[2] --> 788
# mult.multiplier[0] --> 904 908
# mult.multiplier[1] --> 808 812
# mult.multiplier[2] --> 595
# mult.product[0] --> 1000
# mult.product[1] --> 1025
# mult.product[2] --> 782
# mult.product[3] --> 422
# mult.product[4] --> 177
# mult.product[5] --> 367
Claim #1:  3 *  5 = 15 [YES] -- 3 @ -389.75
```
In the above, `[YES]` implies that the multiplication is correct, not that it is necessarily the inputs/outputs the user requested.

One can pin the output and one of the inputs to perform integer division, say 15 ÷ 3 (binary 001111 ÷ 011):
```bash
$ edif2qmasm mult.edif | qmasm --run --all-solns --pin="mult.multiplicand[2:0] := 011" --pin="mult.product[5:0] := 001111" | ./friendly-mult | grep YES
# mult.multiplicand[0] --> 514 519
# mult.multiplicand[1] --> 538 541
# mult.multiplicand[2] --> 627 631
# mult.multiplier[0] --> 521 526
# mult.multiplier[1] --> 880 884
# mult.multiplier[2] --> 298 394
# mult.product[0] --> 705
# mult.product[1] --> 131
# mult.product[2] --> 314
# mult.product[3] --> 165
# mult.product[4] --> 214
# mult.product[5] --> 489
Claim #4:  3 *  5 = 15 [YES] -- 1 @ -435.25
```
Note that we specified the `--all-solns` option to `qmasm` as a workaround until we determine why the minimal-energy solutions are so often incorrect.

Finally, one can pin only the output to factor a number, say 15 (binary 001111):
```bash
$ edif2qmasm mult.edif | qmasm --run --all-solns --pin="mult.product[5:0] := 001111" | ./friendly-mult | grep YES
# mult.multiplicand[0] --> 73
# mult.multiplicand[1] --> 298 300
# mult.multiplicand[2] --> 209 213
# mult.multiplier[0] --> 137 140
# mult.multiplier[1] --> 179 182
# mult.multiplier[2] --> 507
# mult.product[0] --> 37
# mult.product[1] --> 220
# mult.product[2] --> 720 816
# mult.product[3] --> 169 172
# mult.product[4] --> 493
# mult.product[5] --> 363
Claim #21:  3 *  5 = 15 [YES] -- 1 @ -457.75
Claim #39:  5 *  3 = 15 [YES] -- 1 @ -458.75
```
Additional examples
-------------------

More examples are forthcoming, but our first priority is to diagnose why the *mult* example is not behaving as expected.
