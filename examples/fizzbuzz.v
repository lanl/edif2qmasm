// Play the Fizz Buzz game.
//
// Given a number n, set fizz to 1 iff n is divisible by 3, and set
// buzz to 1 iff n is divisible by 5.
//
// Author: Scott Pakin <pakin@lanl.gov>

module fizzbuzz (n, fizz, buzz);
   parameter BITS = 7;
   input [BITS - 1 : 0] n;
   output fizz;
   output buzz;
   wire [(1<<BITS) - 1 : 0] mod3;
   wire [(1<<BITS) - 1 : 0] mod5;

   generate
      genvar i;
      for (i = 0; i < (1<<BITS); i = i + 1)
        assign mod3[i] = i%3 == 0;
      for (i = 0; i < (1<<BITS); i = i + 1)
        assign mod5[i] = i%5 == 0;
   endgenerate

   assign fizz = mod3[n];
   assign buzz = mod5[n];
endmodule // fizzbuzz
