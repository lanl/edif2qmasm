###########################################
# Build edif2qmasm, a converter from EDIF #
# netlists to LANL's QMASM representation #
#					  #
# By Scott Pakin <pakin@lanl.gov>         #
###########################################

# Modify the following as needed.
prefix = /usr/local
bindir = $(prefix)/bin
mandir = $(prefix)/share/man/man1
sharedir = $(prefix)/share/edif2qmasm
GO = go
INSTALL = install
SED = sed

all: edif2qmasm stdcell.qmasm edif2qmasm.1

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
	pigeon parse-edif.peg > parse-edif.tmp
	goimports parse-edif.tmp | gofmt > parse-edif.go
	$(RM) parse-edif.tmp

sexptype_string.go: $(SOURCES_NO_SEXPTYPE)
	stringer -type=SExpType $(SOURCES_NO_SEXPTYPE)

edif2qmasm.1: edif2qmasm.rst
	$(SED) "s/:Date:.*/:Date: $$(date +'%Y-%m-%d')/" edif2qmasm.rst | \
	  rst2man > edif2qmasm.1

clean:
	$(RM) edif2qmasm
	$(RM) parse-edif.tmp

maintainer-clean:
	$(RM) parse-edif.go sexptype_string.go edif2qmasm.1

install: edif2qmasm stdcell.qmasm edif2qmasm.1
	$(INSTALL) -m 0755 -d $(DESTDIR)$(bindir)
	$(INSTALL) -m 0755 edif2qmasm $(DESTDIR)$(bindir)
	$(INSTALL) -m 0755 -d $(DESTDIR)$(sharedir)
	$(INSTALL) -m 0644 stdcell.qmasm $(DESTDIR)$(sharedir)
	$(INSTALL) -m 0755 -d $(DESTDIR)$(mandir)
	$(INSTALL) -m 0644 edif2qmasm.1 $(DESTDIR)$(mandir)
	gzip $(DESTDIR)$(mandir)/edif2qmasm.1

.PHONY: all clean maintainer-clean install
