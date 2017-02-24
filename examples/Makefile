######################################
# Construct various test cases using #
# the Yosys Open SYnthesis Suite     #
#                                    #
# By Scott Pakin <pakin@lanl.gov>    #
######################################

VSOURCES = \
	circsat.v \
	map-color.v \
	mult.v

EDIFS = $(patsubst %.v, %.edif, $(VSOURCES))
QMASMS = $(patsubst %.v, %.qmasm, $(VSOURCES))
YOSYS = yosys
E2Q = ../edif2qmasm

all: edif

edif: $(EDIFS)

%.edif: %.v synth.ys
	$(YOSYS) -q $< synth.ys -b edif -o $@

qmasm: $(QMASMS)

%.qmasm: %.edif
	$(E2Q) -o $@ $<

clean:
	$(RM) $(QMASMS)
	$(RM) $(EDIFS)

.PHONY: all clean edif qmasm
