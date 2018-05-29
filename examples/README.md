edif2qmasm Examples
===================

This directory contains a few examples written in [Verilog](https://en.wikipedia.org/wiki/Verilog) and a `Makefile` that compiles them to [EDIF](https://en.wikipedia.org/wiki/EDIF) netlists using the [Yosys Open SYnthesis Suite](http://www.clifford.at/yosys/).  Start by running
```bash
make
```
to produce a `.edif` file from each `.v` file or use your favorite hardware-synthesis tool to perform the equivalent operation.

`edif2qmasm` supports only the following gates (defined in [stdcell.qmasm](https://github.com/lanl/edif2qmasm/blob/master/stdcell.qmasm)) so all designs must be compiled to use only these:

* 1-input: NOT, DFF_P, DFF_N
* 2-input: AND, NAND, OR, NOR, XOR, XNOR
* 3-input: MUX, AOI3, OIA3
* 4-input: AOI4, OIA4

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

Maximum cut
-----------

The maximum-cut, or max-cut,  decision problem is:

> Given a graph *G* and an integer *k*, is there is a cut of size at least *k* in *G*?

[`max-cut.v`](https://github.com/lanl/edif2qmasm/blob/master/examples/max-cut.v) solves this for the following, hard-wired graph (taken from https://en.wikipedia.org/wiki/Maximum_cut):

![Maximum-cut example](https://upload.wikimedia.org/wikipedia/commons/c/cf/Max-cut.svg)

We label the vertices *A* through *E* in top-to-bottom, left-to-right order.  As the figure shows, the maximum cut of 5 is found when *A* and *E* are colored white and *B*, *C*, and *D* are colored black.  As with many of the examples on this page, the code is run backward from `valid` being `True` to a Boolean coloring of `a` through `e`:
```bash
$ edif2qmasm max-cut.edif | qmasm --run --values=ints --pin="maxcut.valid := true"
# maxcut.a --> 525
# maxcut.b --> 901
# maxcut.c --> 635
# maxcut.cut[0] --> 33
# maxcut.cut[1] --> 541
# maxcut.cut[2] --> 130 133 141
# maxcut.d --> 323 419 515
# maxcut.e --> 366
# maxcut.valid --> 1079
Solution #1 (energy = -494.58, tally = 1):

    Name          Binary  Decimal
    ------------  ------  -------
    maxcut.a           1        1
    maxcut.b           0        0
    maxcut.c           0        0
    maxcut.cut       010        2
    maxcut.d           0        0
    maxcut.e           0        0
    maxcut.valid       1        1

Solution #2 (energy = -494.58, tally = 1):

    Name          Binary  Decimal
    ------------  ------  -------
    maxcut.a           0        0
    maxcut.b           1        1
    maxcut.c           0        0
    maxcut.cut       010        2
    maxcut.d           0        0
    maxcut.e           0        0
    maxcut.valid       1        1
```

The `friendly-maxcut` script post-processes QMASM's output into a more human-readable form:
```bash
$ edif2qmasm max-cut.edif | qmasm --run --values=ints --pin="maxcut.valid := true" | ./friendly-maxcut
# maxcut.a --> 525
# maxcut.b --> 901
# maxcut.c --> 635
# maxcut.cut[0] --> 33
# maxcut.cut[1] --> 541
# maxcut.cut[2] --> 130 133 141
# maxcut.d --> 323 419 515
# maxcut.e --> 366
# maxcut.valid --> 1079
Claim #1: | C | A B D E | 2 >= 2 with tally = 1 and energy = -495.58 [YES]
Claim #2: | C | A B D E | 2 >= 3 with tally = 1 and energy = -495.58 [NO]
```

More interesting usage is to specify the minimum cut size, `cut`.  In our sample graph, the maximum cut one can make is 5 (binary 101):
```bash
$ edif2qmasm max-cut.edif | qmasm --run --values=ints --pin="maxcut.valid := true" --pin="maxcut.cut[2:0] := 101" | ./friendly-maxcut
# maxcut.a --> 387 483 579 583 591
# maxcut.b --> 936 940
# maxcut.c --> 290 294
# maxcut.cut[0] --> 747
# maxcut.cut[1] --> 435
# maxcut.cut[2] --> 597 600 605
# maxcut.d --> 211
# maxcut.e --> 374
# maxcut.valid --> 906
Claim #1: | A | B C D E | 3 >= 5 with tally = 1 and energy = -465.58 [NO]
Claim #2: | A C | B D E | 3 >= 5 with tally = 1 and energy = -465.58 [NO]
```

You may need to use either `--all-solns` or `--postproc=opt` to increase the likelihood of receiving a correct solution:
```bash
$ edif2qmasm max-cut.edif | qmasm --run --values=ints --pin="maxcut.valid := true" --pin="maxcut.cut[2:0] := 101" --all-solns | ./friendly-maxcut
# maxcut.a --> 387 483 579 583 591
# maxcut.b --> 936 940
# maxcut.c --> 290 294
# maxcut.cut[0] --> 747
# maxcut.cut[1] --> 435
# maxcut.cut[2] --> 597 600 605
# maxcut.d --> 211
# maxcut.e --> 374
# maxcut.valid --> 906
Claim #1: | C | A B D E | 2 >= 5 with tally = 1 and energy = -466.08 [NO]
Claim #2: | B | A C D E | 2 >= 5 with tally = 1 and energy = -465.42 [NO]
Claim #3: | | A B C D E | 0 >= 5 with tally = 2 and energy = -465.42 [NO]
Claim #4: | B C D | A E | 5 >= 5 with tally = 1 and energy = -465.08 [YES]
Claim #5: | D | A B C E | 3 >= 5 with tally = 1 and energy = -465.08 [NO]
Claim #6: | A D | B C E | 4 >= 5 with tally = 1 and energy = -465.08 [NO]
Claim #7: | A C | B D E | 3 >= 5 with tally = 1 and energy = -464.58 [NO]
Claim #8: | C D | A B E | 3 >= 5 with tally = 1 and energy = -464.58 [NO]
Claim #9: | A B | C D E | 3 >= 5 with tally = 1 and energy = -464.08 [NO]
Claim #10: | B D | A C E | 5 >= 5 with tally = 1 and energy = -464.08 [YES]
Claim #11: | D E | A B C | 3 >= 5 with tally = 1 and energy = -464.08 [NO]
Claim #12: | A | B C D E | 3 >= 5 with tally = 1 and energy = -464.08 [NO]
Claim #13: | E | A B C D | 2 >= 5 with tally = 1 and energy = -463.58 [NO]
Claim #14: | A B C | D E | 3 >= 5 with tally = 1 and energy = -463.08 [NO]
Claim #15: | A C D | B E | 2 >= 5 with tally = 1 and energy = -461.08 [NO]
$ edif2qmasm max-cut.edif | qmasm --run --values=ints --pin="maxcut.valid := true" --pin="maxcut.cut[2:0] := 101" --postproc=opt | ./friendly-maxcut
# maxcut.a --> 387 483 579 583 591
# maxcut.b --> 936 940
# maxcut.c --> 290 294
# maxcut.cut[0] --> 747
# maxcut.cut[1] --> 435
# maxcut.cut[2] --> 597 600 605
# maxcut.d --> 211
# maxcut.e --> 374
# maxcut.valid --> 906
Claim #1: | B C D | A E | 5 >= 5 with tally = 4 and energy = -467.58 [YES]
Claim #2: | A E | B C D | 5 >= 5 with tally = 7 and energy = -467.58 [YES]
Claim #3: | B D | A C E | 5 >= 5 with tally = 9 and energy = -467.58 [YES]
Claim #4: | A C E | B D | 5 >= 5 with tally = 8 and energy = -467.58 [YES]
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

Fizz Buzz
---------

In [`fizzbuzz.v`](https://github.com/lanl/edif2qmasm/blob/master/examples/fizzbuzz.v) we present, primarily for entertainment value, an implementation of "Quantum FizzBuzz".  [Fizz Buzz](https://en.wikipedia.org/wiki/Fizz_buzz) is a children's game in which players, seated in a circle, count off sequentially.  However, in place of each number that is a multiple of three a player must say "fizz"; in place of each number that is a multiple of five a player must say "buzz"; and in place of each number that is a multiple of both three and five a player must say "fizz buzz".  That is, the first twenty numbers are to be called out as follows:

> 1, 2, fizz, 4, buzz, fizz, 7, 8, fizz, buzz, 11, fizz, 13, 14, fizz buzz, 16, 17, fizz, 19, buzz

The code inputs a 7-bit number `fizzbuzz.n` and reports in outputs `fizzbuzz.fizz` and `fizzbuzz.buzz` whether that number is a multiple of 3, 5, both, or neither.  For example, the number 78 (binary 1001110) is a multiple of 3 but not 5:

```bash
$ edif2qmasm fizzbuzz.edif | qmasm --run -O1 --postproc=opt --values=ints --pin="fizzbuzz.n[6:0] := 1001110"
# fizzbuzz.buzz --> 852
# fizzbuzz.fizz --> 659
# fizzbuzz.n[0] --> 473 478 560 564 569 572 656 752 848
# fizzbuzz.n[1] --> 118 122 126 127 134 142 150 158 166 170 174 179 182 275 371
# fizzbuzz.n[2] --> 56 152 156 157 248 252 344 440 536
# fizzbuzz.n[3] --> 114 210 302 305 306 308 310 316 324 332 336 340 401 402 432 439 447 448 455 498 544 594 640 736 741 832
# fizzbuzz.n[4] --> 387 483 579
# fizzbuzz.n[5] --> 920 924 1004 1012 1016 1020
# fizzbuzz.n[6] --> 1008 1013 1015
Solution #1 (energy = -124.13, tally = 2):

    Name           Binary   Decimal
    -------------  -------  -------
    fizzbuzz.buzz        0        0
    fizzbuzz.fizz        1        1
    fizzbuzz.n     1001110       78
```

Note that the problem is sufficiently large that `--postproc=opt` was needed to nudge the results returned by the hardware towards a correct solution.

Implementing Fizz Buzz for, say, the numbers 1 through 100, is a question commonly asked at interviews for programmer positions.  Given enough samples, [`fizzbuzz.v`](https://github.com/lanl/edif2qmasm/blob/master/examples/fizzbuzz.v) should report all numbers in the range [0, 127] (not likely in order, however):

```bash
$ edif2qmasm fizzbuzz.edif | qmasm --run -O1 --postproc=opt --values=ints --samples=100000
# fizzbuzz.buzz --> 272
# fizzbuzz.fizz --> 101
# fizzbuzz.n[0] --> 5 13 21 29 37 45 53 61 69 77 85
# fizzbuzz.n[1] --> 105 109 120 201 208 211 214 216 222 230 238 297 301 304 309
# fizzbuzz.n[2] --> 328 405 413 419 421 424 425 429 437 515 521 617 713 809 905 1001
# fizzbuzz.n[3] --> 168 245 253 261 264 269 360 456 486 489 494 497 501 502 509 517 525 533 541 549 552 557 585 648 744 840 844 852 933 936 941 1032 1110 1118 1126 1128 1134
# fizzbuzz.n[4] --> 889 891 892
# fizzbuzz.n[5] --> 745 750 835 839 841 847
# fizzbuzz.n[6] --> 748 752 756 848 854
Solution #1 (energy = -124.63, tally = 56):

    Name           Binary   Decimal
    -------------  -------  -------
    fizzbuzz.buzz        0        0
    fizzbuzz.fizz        0        0
    fizzbuzz.n     0000100        4

Solution #2 (energy = -124.63, tally = 67):

    Name           Binary   Decimal
    -------------  -------  -------
    fizzbuzz.buzz        1        1
    fizzbuzz.fizz        0        0
    fizzbuzz.n     0000101        5
```
<span style="text-align: center">…</span>
```bash
Solution #128 (energy = -124.63, tally = 71):

    Name           Binary   Decimal
    -------------  -------  -------
    fizzbuzz.buzz        0        0
    fizzbuzz.fizz        0        0
    fizzbuzz.n     0111110       62
```

Unlike a typical Fizz Buzz implementation, ours can also be used, for example, to output only fizz-buzz numbers:

```bash
$ edif2qmasm fizzbuzz.edif | qmasm --run -O1 --postproc=opt --values=ints --samples=10000 --pin="fizzbuzz.fizz := true" --pin="fizzbuzz.buzz := true" | grep fizzbuzz.n | sort
# fizzbuzz.buzz --> 635
# fizzbuzz.fizz --> 205
# fizzbuzz.n[0] --> 208 212 304 400 496 536 592 597 605 613 621 629 632 633 637
# fizzbuzz.n[1] --> 64 68 160 246 254 256 262 270 278 352 358 359
# fizzbuzz.n[2] --> 147 221 229 237 243 245 339 435 436 437 445 448 453 531 532 533 534 544 627
# fizzbuzz.n[3] --> 33 39 47 55 63 71 79 129 225 321 327 401 409 417 495 497 503 505 511 513 519 527 535 543 551 559 609 705 801 897 901
# fizzbuzz.n[4] --> 911
# fizzbuzz.n[5] --> 841 845 853 861 937
# fizzbuzz.n[6] --> 842 844
    fizzbuzz.n     0000000        0
    fizzbuzz.n     0001111       15
    fizzbuzz.n     0011110       30
    fizzbuzz.n     0101101       45
    fizzbuzz.n     0111100       60
    fizzbuzz.n     1001011       75
    fizzbuzz.n     1011010       90
    fizzbuzz.n     1101001      105
    fizzbuzz.n     1111000      120
```

or, for a particularly contrived example, fizz numbers in which both bits 2 and 5 are True:

```bash
$ edif2qmasm fizzbuzz.edif | qmasm --run -O1 --postproc=opt --values=ints --samples=10000 --pin="fizzbuzz.fizz fizzbuzz.n[5] fizzbuzz.n[2] := 111" | grep fizzbuzz.n | sort
# fizzbuzz.buzz --> 830
# fizzbuzz.fizz --> 979 1075
# fizzbuzz.n[0] --> 816 888 894 902 910 912 918 926 984 1078 1080 1086
# fizzbuzz.n[1] --> 300 304 308 314 316 410 412 506 602 606 698 794 890 986
# fizzbuzz.n[2] --> 715 811 813 907 1003 1004 1005 1012 1020
# fizzbuzz.n[3] --> 136 140 141 232 235 236 244 331 427 523 591 592 599 607 615 619 623 631 639 642 647 738 834 838 930
# fizzbuzz.n[4] --> 434 437 530 626 629
# fizzbuzz.n[5] --> 58 154
# fizzbuzz.n[6] --> 156
    fizzbuzz.n     0100100       36
    fizzbuzz.n     0100111       39
    fizzbuzz.n     0101101       45
    fizzbuzz.n     0110110       54
    fizzbuzz.n     0111100       60
    fizzbuzz.n     0111111       63
    fizzbuzz.n     1100110      102
    fizzbuzz.n     1101100      108
    fizzbuzz.n     1101111      111
    fizzbuzz.n     1110101      117
    fizzbuzz.n     1111110      126
```
