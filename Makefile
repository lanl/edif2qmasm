###########################################
# Build edif2qmasm, a converter from EDIF #
# netlists to LANL's QMASM representation #
#					  #
# By Scott Pakin <pakin@lanl.gov>         #
###########################################

GO = go
PIGEON = pigeon

all: edif2qmasm

SOURCES = \
	edif.go \
	edif2qmasm.go \
	parse-edif.go \
	qmasm.go \
	sexptype_string.go \
	walk-sexp.go

SOURCES_NO_SEXPTYPE = $(subst sexptype_string.go,,$(SOURCES))

edif2qmasm: $(SOURCES)
	$(GO) build -o edif2qmasm $(SOURCES)

parse-edif.go: parse-edif.peg
	$(PIGEON) parse-edif.peg > parse-edif.tmp
	goimports parse-edif.tmp | gofmt > parse-edif.go
	$(RM) parse-edif.tmp

sexptype_string.go: $(SOURCES_NO_SEXPTYPE)
	stringer -type=SExpType $(SOURCES_NO_SEXPTYPE)

vet: $(SOURCES)
	$(GO) vet $(SOURCES)

clean:
	$(RM) edif2qmasm
	$(RM) parse-edif.tmp parse-edif.go
	$(RM) sexptype_string.go

.PHONY: all clean vet
