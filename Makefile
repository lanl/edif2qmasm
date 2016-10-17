##########################################
# Build edif2qasm, a converter from EDIF #
# netlists to LANL's QASM representation #
#					 #
# By Scott Pakin <pakin@lanl.gov>        #
##########################################

AWK = awk
GO = go
PIGEON = pigeon

all: edif2qubo

EDIF2QUBO_SOURCES = \
	edif2qubo.go \
	parse-edif.go \
	build-graph.go \
	netlist-types.go \
	qubos.go \
	sexptype_string.go \
	quboelttype_string.go

edif2qubo: $(EDIF2QUBO_SOURCES)
	$(GO) build -o edif2qubo $(EDIF2QUBO_SOURCES)

parse-edif.go: parse-edif.peg
	$(PIGEON) parse-edif.peg > parse-edif.tmp
	goimports parse-edif.tmp | gofmt > parse-edif.go
	$(RM) parse-edif.tmp

sexptype_string.go: parse-edif.go
	stringer -type=SExpType parse-edif.go

quboelttype_string.go: qubos.go parse-edif.go edif2qubo.go build-graph.go netlist-types.go
	stringer -type=QuboEltType qubos.go parse-edif.go edif2qubo.go build-graph.go netlist-types.go

netlist-types.go: netlist-types.tmpl generate-cells.awk
	cp netlist-types.tmpl netlist-types.go
	$(AWK) -v TTAG=Vcc -v TTEXT=input-power -f generate-cells.awk netlist-types.tmpl >> netlist-types.go
	$(AWK) -v TTAG=Not -v TTEXT=binary-not -f generate-cells.awk netlist-types.tmpl >> netlist-types.go
	$(AWK) -v TTAG=And -v TTEXT=binary-and -f generate-cells.awk netlist-types.tmpl >> netlist-types.go
	$(AWK) -v TTAG=Mux -v TTEXT=multiplexer -f generate-cells.awk netlist-types.tmpl >> netlist-types.go
	$(AWK) -v TTAG=User -v TTEXT=user-defined -f generate-cells.awk netlist-types.tmpl >> netlist-types.go

vet-edif2qubo: $(EDIF2QUBO_SOURCES)
	$(GO) vet $(EDIF2QUBO_SOURCES)

clean:
	$(RM) edif2qubo
	$(RM) parse-edif.tmp parse-edif.go
	$(RM) netlist-types.go sexptype_string.go quboelttype_string.go 

.PHONY: all clean vet-edif2qubo
