##########################################
# Build edif2qasm, a converter from EDIF #
# netlists to LANL's QASM representation #
#					 #
# By Scott Pakin <pakin@lanl.gov>        #
##########################################

GO = go
PIGEON = pigeon

all: edif2qasm

SOURCES = \
	edif.go \
	edif2qasm.go \
	parse-edif.go \
	qasm.go \
	sexptype_string.go \
	walk-sexp.go

SOURCES_NO_SEXPTYPE = $(subst sexptype_string.go,,$(SOURCES))

edif2qasm: $(SOURCES)
	$(GO) build -o edif2qasm $(SOURCES)

parse-edif.go: parse-edif.peg
	$(PIGEON) parse-edif.peg > parse-edif.tmp
	goimports parse-edif.tmp | gofmt > parse-edif.go
	$(RM) parse-edif.tmp

sexptype_string.go: $(SOURCES_NO_SEXPTYPE)
	stringer -type=SExpType $(SOURCES_NO_SEXPTYPE)

vet: $(SOURCES)
	$(GO) vet $(SOURCES)

clean:
	$(RM) edif2qasm
	$(RM) parse-edif.tmp parse-edif.go
	$(RM) sexptype_string.go

.PHONY: all clean vet
