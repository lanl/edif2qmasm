// Multiply a 4-bit multiplicand by a 4-bit multiplier to get an 8-bit
// product.  Specify the product and either the multiplicand or
// multiplier to perform integer division.  Specify only the product
// to factor it into two integers.
//
// Author: Scott Pakin <pakin@lanl.gov>

module mult (multiplicand, multiplier, product);
   input [3:0] multiplicand;
   input [3:0] multiplier;
   output[7:0] product;

   assign product = multiplicand * multiplier;
endmodule // mult
