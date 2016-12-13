// Say whether a given four-color map coloring is valid.
// The map corresponds to the regions in the Land of Oz:
//
//                     [Gillikin Country]
//                    /         |        \
// [Munchkin Country] -- [Emerald City] -- [Winkie Country]
//                    \         |        /
//                     [Quadling Country]
//
// Author: Scott Pakin <pakin@lanl.gov>

module map_color (GC, WC, QC, MC, EC, valid);
   input [1:0] GC;
   input [1:0] WC;
   input [1:0] QC;
   input [1:0] MC;
   input [1:0] EC;
   output      valid;
   wire [7:0]  tests;

   assign tests[0] = GC != WC;
   assign tests[1] = WC != QC;
   assign tests[2] = QC != MC;
   assign tests[3] = MC != GC;
   assign tests[4] = EC != GC;
   assign tests[5] = EC != WC;
   assign tests[6] = EC != QC;
   assign tests[7] = EC != MC;

   assign valid = &tests[7:0];
endmodule // map_color
