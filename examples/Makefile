######################################
# Construct various test cases using #
# the Yosys Open SYnthesis Suite     #
#                                    #
# By Scott Pakin <pakin@lanl.gov>    #
######################################

VSOURCES = \
	circsat.v

EDIFS = $(patsubst %.v, %.edif, $(VSOURCES))
YOSYS = yosys

all: $(EDIFS)

%.edif: %.v synth.ys
	$(YOSYS) -q $< synth.ys -b edif -o $@

clean:
	$(RM) $(EDIFS)

.PHONY: all clean
