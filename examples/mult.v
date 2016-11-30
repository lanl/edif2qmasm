// Multiply a 3-bit multiplicand by a 3-bit multiplier to get a 6-bit
// product.  Specify the product and either the multiplicand or
// multiplier to perform integer division.  Specify only the product
// to factor it into two integers.
//
// Author: Scott Pakin <pakin@lanl.gov>

module mult (multiplicand, multiplier, product);
   input [2:0] multiplicand;
   input [2:0] multiplier;
   output[5:0] product;

   assign product = multiplicand * multiplier;
endmodule // mult
