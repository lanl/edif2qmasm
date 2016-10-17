##########################################
# Build edif2qasm, a converter from EDIF #
# netlists to LANL's QASM representation #
#					 #
# By Scott Pakin <pakin@lanl.gov>        #
##########################################

GO = go
PIGEON = pigeon

all: edif2qasm

EDIF2QUBO_SOURCES = \
	edif2qasm.go \
	parse-edif.go \
	sexptype_string.go

edif2qasm: $(EDIF2QUBO_SOURCES)
	$(GO) build -o edif2qasm $(EDIF2QUBO_SOURCES)

parse-edif.go: parse-edif.peg
	$(PIGEON) parse-edif.peg > parse-edif.tmp
	goimports parse-edif.tmp | gofmt > parse-edif.go
	$(RM) parse-edif.tmp

sexptype_string.go: parse-edif.go
	stringer -type=SExpType parse-edif.go

vet-edif2qasm: $(EDIF2QUBO_SOURCES)
	$(GO) vet $(EDIF2QUBO_SOURCES)

clean:
	$(RM) edif2qasm
	$(RM) parse-edif.tmp parse-edif.go
	$(RM) sexptype_string.go

.PHONY: all clean vet-edif2qasm
