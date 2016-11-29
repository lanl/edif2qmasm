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

VERSION = 1.0

GEN_SOURCES = \
	parse-edif.go \
	sexptype_string.go

REG_SOURCES = \
	edif.go \
	edif2qmasm.go \
	qmasm.go \
	walk-sexp.go

SOURCES = $(REG_SOURCES) $(GEN_SOURCES)

TARCONTENTS = \
	$(SOURCES) \
	edif2qmasm.1 \
	edif2qmasm.rst \
	Makefile \
	parse-edif.peg \
	stdcell.qmasm

edif2qmasm: $(SOURCES)
	$(GO) build -o edif2qmasm

parse-edif.go: parse-edif.peg $(REG_SOURCES)
	$(GO) generate -x

# The rule for parse-edif.go produces sexptype_string.go as a side effect.
sexptype_string.go: parse-edif.go

edif2qmasm.1: edif2qmasm.rst
	$(SED) "s/:Date:.*/:Date: $$(date +'%Y-%m-%d')/" edif2qmasm.rst | \
	  rst2man > edif2qmasm.1

clean:
	$(RM) edif2qmasm

maintainer-clean:
	$(RM) $(GEN_SOURCES) parse-edif.tmp
	$(RM) edif2qmasm.1 edif2qmasm-$(VERSION).tar.gz

install: edif2qmasm stdcell.qmasm edif2qmasm.1
	$(INSTALL) -m 0755 -d $(DESTDIR)$(bindir)
	$(INSTALL) -m 0755 edif2qmasm $(DESTDIR)$(bindir)
	$(INSTALL) -m 0755 -d $(DESTDIR)$(sharedir)
	$(INSTALL) -m 0644 stdcell.qmasm $(DESTDIR)$(sharedir)
	$(INSTALL) -m 0755 -d $(DESTDIR)$(mandir)
	$(INSTALL) -m 0644 edif2qmasm.1 $(DESTDIR)$(mandir)
	gzip $(DESTDIR)$(mandir)/edif2qmasm.1

dist: edif2qmasm-$(VERSION).tar.gz

edif2qmasm-$(VERSION).tar.gz: $(TARCONTENTS)
	$(RM) -r edif2qmasm-$(VERSION)
	mkdir edif2qmasm-$(VERSION)
	cp $(TARCONTENTS) edif2qmasm-$(VERSION)
	tar -czvf edif2qmasm-$(VERSION).tar.gz edif2qmasm-$(VERSION)
	$(RM) -r edif2qmasm-$(VERSION)

.PHONY: all clean maintainer-clean install dist
