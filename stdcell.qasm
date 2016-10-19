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
