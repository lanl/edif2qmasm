edif2qmasm Examples
===================

This directory contains a few examples written in [Verilog](https://en.wikipedia.org/wiki/Verilog) and a `Makefile` that compiles them to [EDIF](https://en.wikipedia.org/wiki/EDIF) netlists using the [Yosys Open SYnthesis Suite](http://www.clifford.at/yosys/).  Start by running
```bash
make
```
to produce a `.edif` file from each `.v` file or use your favorite hardware-synthesis tool to perform the equivalent operation.

`edif2qmasm` currently supports only a handful of gates (defined in [stdcell.qmasm](https://github.com/lanl/edif2qmasm/blob/master/stdcell.qmasm)) so all designs must be compiled to use only those gates.

* 1-input: NOT
* 2-input: AND, OR, XOR
* 3-input: MUX

The rest of this document describes each of the examples in turn.

circsat
-------

[`circsat.v`](https://github.com/lanl/edif2qmasm/blob/master/examples/circsat.v) implements an arbitrary circuit, specifically the one presented in [Cormen, Leiserson, Rivest, and Stein's Algorithms textbook](https://mitpress.mit.edu/books/introduction-algorithms) in its discussion of circuit-satisfiability problems.  The Verilog module's inputs are named *a*, *b*, and *c*, and the sole output is named *y*.  All are single bits.  Internally, inputs, outputs, and intermediate values are named *x*[1]…*x*[10].

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

Excluding the comments, module definition, and parameter declarations, [`mult.v`](https://github.com/lanl/edif2qmasm/blob/master/examples/mult.v) is a one-line Verilog program:
```Verilog
assign product = multiplicand * multiplier;
```

`mult.v` showcases the flexibility of `edif2qmasm`'s approach.  For starters, one can pin values to the inputs, say 3 and 5 (binary 011 and 101) to multiply them together:
```bash
$ edif2qmasm mult.edif | qmasm --run --pin="mult.multiplicand[2:0] := 011" --pin="mult.multiplier[2:0] := 101"
# mult.multiplicand[0] --> 810 815 823 831 906
# mult.multiplicand[1] --> 1016
# mult.multiplicand[2] --> 708
# mult.multiplier[0] --> 329 335
# mult.multiplier[1] --> 268
# mult.multiplier[2] --> 690
# mult.product[0] --> 445
# mult.product[1] --> 614
# mult.product[2] --> 716
# mult.product[3] --> 627
# mult.product[4] --> 461
# mult.product[5] --> 1108
Solution #1 (energy = -497.75, tally = 489):

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
    mult.product[3]            +1  True
    mult.product[4]            -1  False
    mult.product[5]            -1  False
```

The `friendly-mult` script post-processes the above into a more human-readable form:
```bash
$ edif2qmasm mult.edif | qmasm --run --pin="mult.multiplicand[2:0] := 011" --pin="mult.multiplier[2:0] := 101" | ./friendly-mult
# mult.multiplicand[0] --> 300
# mult.multiplicand[1] --> 412 420
# mult.multiplicand[2] --> 887
# mult.multiplier[0] --> 213 221
# mult.multiplier[1] --> 970 975
# mult.multiplier[2] --> 118
# mult.product[0] --> 198
# mult.product[1] --> 778
# mult.product[2] --> 246
# mult.product[3] --> 804
# mult.product[4] --> 624
# mult.product[5] --> 838
Claim #1:  3 *  5 = 15 [YES] -- 450 @ -444.75
```
In the above, `[YES]` implies that the multiplication is correct, not that it is necessarily the inputs/outputs the user requested.

One can pin the output and one of the inputs to perform integer division, say 15 ÷ 3 (binary 001111 ÷ 011):
```bash
$ edif2qmasm mult.edif | qmasm --run --pin="mult.multiplicand[2:0] := 011" --pin="mult.product[5:0] := 001111" | ./friendly-mult
# mult.multiplicand[0] --> 251 347 351
# mult.multiplicand[1] --> 852 859 860
# mult.multiplicand[2] --> 748
# mult.multiplier[0] --> 368
# mult.multiplier[1] --> 729 734
# mult.multiplier[2] --> 639
# mult.product[0] --> 332
# mult.product[1] --> 874
# mult.product[2] --> 426
# mult.product[3] --> 499
# mult.product[4] --> 454 457 462
# mult.product[5] --> 933
Claim #1:  3 *  5 = 15 [YES] -- 103 @ -474.75
```

Finally, one can pin only the output to factor a number, say 15 (binary 001111):
```bash
$ edif2qmasm mult.edif | qmasm --run --pin="mult.product[5:0] := 001111" | ./friendly-mult
# mult.multiplicand[0] --> 545 550
# mult.multiplicand[1] --> 240
# mult.multiplicand[2] --> 946 950
# mult.multiplier[0] --> 723 726
# mult.multiplier[1] --> 357 365 373
# mult.multiplier[2] --> 648 655
# mult.product[0] --> 335
# mult.product[1] --> 421
# mult.product[2] --> 982
# mult.product[3] --> 210
# mult.product[4] --> 819
# mult.product[5] --> 1024
Claim #1:  3 *  5 = 15 [YES] -- 42 @ -460.75
Claim #2:  5 *  3 = 15 [YES] -- 53 @ -460.75
```

Map coloring
------------

Map coloring is arguably the closest thing to a D-Wave "Hello, world!" program.  [`map-color.v`](https://github.com/lanl/edif2qmasm/blob/master/examples/map-color.v) shows how map coloring can be implemented in Verilog with the intention of executing it on a D-Wave system.  The goal of map coloring is to color a given map using at most four colors such that no two adjacent regions use the same color.  `map-color.v` colors a map of the Land of Oz:

![Map of the Land of Oz](https://upload.wikimedia.org/wikipedia/commons/8/8e/Map-of-Oz.jpg)

The interesting aspect of the code is that it is expressed in the *reverse* direction: Given a map coloring, say whether or not it is valid.  By pinning `valid` to *true*, the program returns a set of valid map colorings:
```bash
$ edif2qmasm map_color.edif | qmasm --run --pin="map_color.valid := true"
# map_color.EC[0] --> 163 166 174 259 355
# map_color.EC[1] --> 521
# map_color.GC[0] --> 279
# map_color.GC[1] --> 593 596 604
# map_color.MC[0] --> 80 87
# map_color.MC[1] --> 708 716
# map_color.QC[0] --> 208 304
# map_color.QC[1] --> 665 671
# map_color.WC[0] --> 242
# map_color.WC[1] --> 457 553 649
# map_color.valid --> 741
Solution #1 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       -1  False
    map_color.EC[1]       +1  True
    map_color.GC[0]       +1  True
    map_color.GC[1]       +1  True
    map_color.MC[0]       +1  True
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       +1  True
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #2 (energy = -462.25, tally = 3):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       -1  False
    map_color.EC[1]       +1  True
    map_color.GC[0]       -1  False
    map_color.GC[1]       -1  False
    map_color.MC[0]       +1  True
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       +1  True
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #3 (energy = -462.25, tally = 4):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       +1  True
    map_color.GC[0]       -1  False
    map_color.GC[1]       -1  False
    map_color.MC[0]       +1  True
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       +1  True
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #4 (energy = -462.25, tally = 4):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       +1  True
    map_color.GC[0]       -1  False
    map_color.GC[1]       +1  True
    map_color.MC[0]       +1  True
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       +1  True
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #5 (energy = -462.25, tally = 2):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       +1  True
    map_color.GC[0]       -1  False
    map_color.GC[1]       -1  False
    map_color.MC[0]       -1  False
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       +1  True
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #6 (energy = -462.25, tally = 2):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       -1  False
    map_color.GC[0]       -1  False
    map_color.GC[1]       -1  False
    map_color.MC[0]       -1  False
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       -1  False
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #7 (energy = -462.25, tally = 6):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       -1  False
    map_color.GC[0]       +1  True
    map_color.GC[1]       +1  True
    map_color.MC[0]       -1  False
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       -1  False
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #8 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       -1  False
    map_color.GC[0]       -1  False
    map_color.GC[1]       -1  False
    map_color.MC[0]       +1  True
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       -1  False
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #9 (energy = -462.25, tally = 2):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       +1  True
    map_color.GC[0]       -1  False
    map_color.GC[1]       -1  False
    map_color.MC[0]       +1  True
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       -1  False
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #10 (energy = -462.25, tally = 2):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       +1  True
    map_color.GC[0]       -1  False
    map_color.GC[1]       -1  False
    map_color.MC[0]       -1  False
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       -1  False
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #11 (energy = -462.25, tally = 7):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       -1  False
    map_color.GC[0]       -1  False
    map_color.GC[1]       -1  False
    map_color.MC[0]       -1  False
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       +1  True
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #12 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       -1  False
    map_color.GC[0]       -1  False
    map_color.GC[1]       -1  False
    map_color.MC[0]       +1  True
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       +1  True
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #13 (energy = -462.25, tally = 5):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       -1  False
    map_color.GC[0]       -1  False
    map_color.GC[1]       +1  True
    map_color.MC[0]       +1  True
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       -1  False
    map_color.WC[0]       +1  True
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #14 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       -1  False
    map_color.GC[0]       -1  False
    map_color.GC[1]       +1  True
    map_color.MC[0]       -1  False
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       -1  False
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #15 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       -1  False
    map_color.GC[0]       +1  True
    map_color.GC[1]       +1  True
    map_color.MC[0]       -1  False
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       -1  False
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #16 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       +1  True
    map_color.GC[0]       -1  False
    map_color.GC[1]       +1  True
    map_color.MC[0]       -1  False
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       -1  False
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #17 (energy = -462.25, tally = 4):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       +1  True
    map_color.GC[0]       +1  True
    map_color.GC[1]       -1  False
    map_color.MC[0]       -1  False
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       -1  False
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #18 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       +1  True
    map_color.GC[0]       -1  False
    map_color.GC[1]       +1  True
    map_color.MC[0]       +1  True
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       -1  False
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #19 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       -1  False
    map_color.EC[1]       -1  False
    map_color.GC[0]       -1  False
    map_color.GC[1]       +1  True
    map_color.MC[0]       +1  True
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       +1  True
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #20 (energy = -462.25, tally = 2):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       +1  True
    map_color.GC[0]       -1  False
    map_color.GC[1]       +1  True
    map_color.MC[0]       -1  False
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       +1  True
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #21 (energy = -462.25, tally = 2):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       +1  True
    map_color.GC[0]       -1  False
    map_color.GC[1]       -1  False
    map_color.MC[0]       +1  True
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       +1  True
    map_color.WC[1]       -1  False
    map_color.valid       +1  True

Solution #22 (energy = -462.25, tally = 2):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       -1  False
    map_color.EC[1]       -1  False
    map_color.GC[0]       -1  False
    map_color.GC[1]       +1  True
    map_color.MC[0]       +1  True
    map_color.MC[1]       -1  False
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       +1  True
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #23 (energy = -462.25, tally = 2):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       -1  False
    map_color.EC[1]       -1  False
    map_color.GC[0]       -1  False
    map_color.GC[1]       +1  True
    map_color.MC[0]       +1  True
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       +1  True
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #24 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       -1  False
    map_color.EC[1]       -1  False
    map_color.GC[0]       +1  True
    map_color.GC[1]       -1  False
    map_color.MC[0]       +1  True
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       +1  True
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #25 (energy = -462.25, tally = 3):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       -1  False
    map_color.GC[0]       -1  False
    map_color.GC[1]       -1  False
    map_color.MC[0]       +1  True
    map_color.MC[1]       +1  True
    map_color.QC[0]       -1  False
    map_color.QC[1]       +1  True
    map_color.WC[0]       +1  True
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #26 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       +1  True
    map_color.EC[1]       +1  True
    map_color.GC[0]       +1  True
    map_color.GC[1]       -1  False
    map_color.MC[0]       -1  False
    map_color.MC[1]       +1  True
    map_color.QC[0]       +1  True
    map_color.QC[1]       -1  False
    map_color.WC[0]       -1  False
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #27 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       -1  False
    map_color.EC[1]       -1  False
    map_color.GC[0]       +1  True
    map_color.GC[1]       -1  False
    map_color.MC[0]       -1  False
    map_color.MC[1]       +1  True
    map_color.QC[0]       +1  True
    map_color.QC[1]       -1  False
    map_color.WC[0]       +1  True
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #28 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       -1  False
    map_color.EC[1]       -1  False
    map_color.GC[0]       -1  False
    map_color.GC[1]       +1  True
    map_color.MC[0]       +1  True
    map_color.MC[1]       +1  True
    map_color.QC[0]       +1  True
    map_color.QC[1]       -1  False
    map_color.WC[0]       +1  True
    map_color.WC[1]       +1  True
    map_color.valid       +1  True

Solution #29 (energy = -462.25, tally = 1):

    Name(s)             Spin  Boolean
    ------------------  ----  --------
    map_color.EC[0]       -1  False
    map_color.EC[1]       -1  False
    map_color.GC[0]       +1  True
    map_color.GC[1]       -1  False
    map_color.MC[0]       -1  False
    map_color.MC[1]       +1  True
    map_color.QC[0]       +1  True
    map_color.QC[1]       +1  True
    map_color.WC[0]       -1  False
    map_color.WC[1]       +1  True
    map_color.valid       +1  True
```

The `friendly-map` script post-processes the above into a more human-readable form:
```bash
$ edif2qmasm map-color.edif | qmasm --run --pin="map_color.valid := true" | ./friendly-map
# map_color.EC[0] --> 353 449
# map_color.EC[1] --> 203 205
# map_color.GC[0] --> 536 542 550 558 566
# map_color.GC[1] --> 303
# map_color.MC[0] --> 427 523 527
# map_color.MC[1] --> 602 604 607 615 617 623 713
# map_color.QC[0] --> 77 85
# map_color.QC[1] --> 801
# map_color.WC[0] --> 640
# map_color.WC[1] --> 781 787 789
# map_color.valid --> 108 116 124
Claim #1: EC=3 GC=0 MC=1 QC=0 WC=1 --> True with tally = 1 and energy = -492.25 [YES]
Claim #2: EC=3 GC=0 MC=1 QC=0 WC=2 --> True with tally = 1 and energy = -492.25 [YES]
Claim #3: EC=3 GC=0 MC=2 QC=0 WC=2 --> True with tally = 1 and energy = -492.25 [YES]
Claim #4: EC=1 GC=0 MC=2 QC=0 WC=3 --> True with tally = 1 and energy = -492.25 [YES]
Claim #5: EC=1 GC=2 MC=3 QC=0 WC=3 --> True with tally = 1 and energy = -492.25 [YES]
Claim #6: EC=3 GC=2 MC=0 QC=2 WC=0 --> True with tally = 1 and energy = -492.25 [YES]
Claim #7: EC=3 GC=2 MC=0 QC=2 WC=1 --> True with tally = 3 and energy = -492.25 [YES]
Claim #8: EC=3 GC=0 MC=1 QC=2 WC=1 --> True with tally = 1 and energy = -492.25 [YES]
Claim #9: EC=2 GC=3 MC=0 QC=1 WC=0 --> True with tally = 1 and energy = -492.25 [YES]
Claim #10: EC=3 GC=2 MC=0 QC=1 WC=0 --> True with tally = 1 and energy = -492.25 [YES]
Claim #11: EC=2 GC=1 MC=0 QC=3 WC=0 --> True with tally = 2 and energy = -492.25 [YES]
Claim #12: EC=2 GC=3 MC=0 QC=3 WC=0 --> True with tally = 3 and energy = -492.25 [YES]
Claim #13: EC=2 GC=3 MC=0 QC=3 WC=1 --> True with tally = 1 and energy = -492.25 [YES]
```
