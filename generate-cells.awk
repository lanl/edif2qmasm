# This file is needed to build edif2qubo.  It generates cells and cell
# instances from templates defined in netlist-types.tmpl.
#
# Usage: awk -v TTAG=Vcc -v TTEXT=input-power -f generate-cells.awk netlist-types.tmpl

BEGIN {
    print ""
}

$0 ~ "BEGIN: Gnd", $0 ~ "END: Gnd" {
    gsub("Gnd", TTAG)
    gsub("ground", TTEXT)
    print
}
