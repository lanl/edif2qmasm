###################################
# QASM standard-cell library	  #
# for use with edif2qasm	  #
#				  #
# By Scott Pakin <pakin@lanl.gov> #
###################################

# N.B. Weights and strengths are currently scaled so that the sum of
# their absolute values is 1.0.  This may change in a future version
# of this library.

# Y = A AND B
!begin_macro AND
A -0.1111
B -0.1111
Y  0.2222

A B  0.1111
A Y -0.2222
B Y -0.2222
!end_macro AND

# Y = A OR B
!begin_macro OR
A  0.1111
B  0.1111
Y -0.2222

A B  0.1111
A Y -0.2222
B Y -0.2222
!end_macro OR

# Y = NOT A
!begin_macro NOT
A Y 1.0
!end_macro NOT

# Y = A XOR B
!begin_macro XOR
A    0.0714
B   -0.0714
Y   -0.0714
$a1 -0.1429

A B   -0.0714
A Y   -0.0714
A $a1 -0.1429
B Y    0.0714
B $a1  0.1429
Y $a1  0.1429
!end_macro XOR

# Constants for power and ground.
!alias VCC true
!alias GND false
