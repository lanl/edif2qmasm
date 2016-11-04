###########################################
# Build edif2qmasm, a converter from EDIF #
# netlists to LANL's QMASM representation #
#					  #
# By Scott Pakin <pakin@lanl.gov>         #
###########################################

# Modify the following as needed.
prefix = /usr/local
bindir = $(prefix)/bin
sharedir = $(prefix)/share/edif2qmasm
GO = go
PIGEON = pigeon
INSTALL = install

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

install: edif2qmasm stdcell.qmasm
	$(INSTALL) -m 0755 -d $(DESTDIR)$(bindir)
	$(INSTALL) -m 0755 edif2qmasm $(DESTDIR)$(bindir)
	$(INSTALL) -m 0755 -d $(DESTDIR)$(sharedir)
	$(INSTALL) -m 0644 stdcell.qmasm $(DESTDIR)$(sharedir)

.PHONY: all clean vet install
