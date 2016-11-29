module circsat (a, b, c, y);
   input a, b, c;
   output y;
   wire [1:10] x;

   assign x[1] = a;
   assign x[2] = b;
   assign x[3] = c;
   assign x[4] = ~x[3];
   assign x[5] = x[1] | x[2];
   assign x[6] = ~x[4];
   assign x[7] = x[1] & x[2] & x[4];
   assign x[8] = x[5] | x[6];
   assign x[9] = x[6] | x[7];
   assign x[10] = x[8] & x[9] & x[7];
   assign y = x[10];
endmodule // circsat
