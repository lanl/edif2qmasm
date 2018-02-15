// Answer a max-cut decision problem.
//
// The max-cut decision problem is, "Given a graph G and an integer k,
// is there is a cut of size at least k in G?"  We solve this for the
// following, hard-wired graph (taken from
// https://en.wikipedia.org/wiki/Maximum_cut):
//
//          A----B
//         /|    |
//        C |    |
//         \|    |
//          D----E
//
// Author: Scott Pakin <pakin@lanl.gov>

`define BITS 3
`define EDGE(X, Y) {(`BITS-1)'b0, (X != Y)}

module maxcut (a, b, c, d, e, cut, valid);
   input             a, b, c, d, e;
   input [`BITS-1:0] cut;
   output            valid;

   assign valid = `EDGE(a, b) + `EDGE(a, c) + `EDGE(a, d) +
                  `EDGE(b, e) + `EDGE(c, d) + `EDGE(d, e) >= cut;
endmodule // maxcut
