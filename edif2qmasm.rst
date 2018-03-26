==========
edif2qmasm
==========

-------------------------------------------
convert EDIF netlists to QMASM source files
-------------------------------------------

:Author: pakin@lanl.gov
:Date: 2018-03-26
:Copyright: BSD
:Version: 1.0
:Manual section: 1

SYNOPSIS
========

    **edif2qmasm** [-o *outfile.qmasm*] [--cycles=\ *N*] *infile.edif*

DESCRIPTION
===========

**ediff2qmasm** converts a hardware circuit specified as an EDIF
netlist to a symbolic Hamiltonian suitable for running on a D-Wave
quantum annealer using QMASM.

Typical usage is to define a circuit using a hardware-description
language (HDL) such as Verilog then passing this to a synthesis tool
like **yosys** to compile the HDL code to EDIF format.  **edif2qmasm**
can then be run on the result, and the resulting QMASM code can be fed
to **qmasm** for execution on a D-Wave system::

    edif2qmasm -o something.qmasm something.edif
    qmasm --run something.qmasm

Optionally, these steps can be combined into a single shell pipeline:

    edif2qmasm something.edif | qmasm --run

OPTIONS
=======

``-o`` *file.qmasm*, ``--output=``\ *file.qmasm*
  Specify the name of the QMASM file to generate.  The default is to
  write QMASM code to the standard output device.

``--cycles=``\ *N*
  Replicate the entire circuit *N* times.  This is used to support
  sequential logic, which needs to be statically unrolled once per
  cycle for QMASM execution.

NOTES
=====

The following is an example of an interactive **yosys** session that
compiles a Verilog source file, *mycircuit.v* to an EDIF netlist,
*mycircuit.edif*::

    yosys> read_verilog mycircuit.v
    yosys> hierarchy; proc; opt; fsm; opt; techmap; opt; clean
    yosys> write_edif mycircuit.edif

For convenience, one might want to create a script, say *synth.ys*,
with contents like the following::

    ###############################################################
    # Generic synthesis script derived from the Yosys README file #
    #                                                             #
    # Usage: yosys infile.v synth.ys -b edif -o outfile.edif      #
    ###############################################################

    # Check design hierarchy.
    hierarchy

    # Translate processes.
    proc; opt

    # Detect and optimize FSM encodings.
    fsm; opt

    # Convert to gate logic.
    techmap; opt

    # Clean up.
    clean

This can then be run conveniently from the command line::

    yosys mycircuit.v synth.ys -b edif -o mycircuit.edif

SEE ALSO
========

yosys(1),
`the QMASM wiki <https://github.com/lanl/qmasm/wiki>`__,
`Wikipedia's entry on Verilog <https://en.wikipedia.org/wiki/Verilog>`__,
`Wikipedia's entry on VHDL <https://en.wikipedia.org/wiki/VHDL>`__,
`Wikipedia's entry on EDIF <https://en.wikipedia.org/wiki/EDIF>`__,
`D-Wave's home page <http://www.dwavesys.com/>`__
