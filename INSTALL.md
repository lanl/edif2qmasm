edif2qmasm installation
=======================

`edif2qmasm` is written in [Go](https://go.dev/) and therefore depends upon a Go compiler to build.

Note that `edif2qmasm` is of limited use withut a compiler than can produce EDIF netlists and [QMASM](https://github.com/lanl/qmasm), which executes the generated code on a D-Wave system.  To date, `edif2qmasm` has been tested only with the [Yosys Open SYnthesis Suite](https://yosyshq.net/yosys/), but reports of usage with other synthesis tools (successful or not) are welcome.

User builds
-----------

There are two ways to build `edif2qmasm`: the `go install` approach and the `make` approach.

### The `go install` approach

Download, build, and install `edif2qmasm` (into your `$GOPATH/bin/` directory) with
```bash
go install github.com/lanl/edif2qmasm@latest
```

You'll also need to copy `stdcell.qmasm` somewhere, say into `/usr/local/share/edif2qmasm/`.  Optionally, you can install the `edif2qmasm.1` man page, say into `/usr/local/share/man/man1/`.

### The `make` approach

As an alternative installation procedure, one can download the code explicitly and build it using the supplied `Makefile`:
```bash
git clone https://github.com/lanl/edif2qmasm.git
cd edif2qmasm
make
make install
```

This approach is supported because it provides a few extra benefits over the simpler, `go install` approach:

* `make clean` (and `make maintainer-clean`) can be used to clean up the build directory.
* `make install` installs the binary, the standard-cell library (`stdcell.qmasm`), and the Unix man page into their standard locations
* `make install` honors the `DESTDIR`, `prefix`, and similar variables, which can be used to override the default installation directories.

Developer builds
----------------

If you want to modify `edif2qmasm` and rebuild it, you'll need a few additional tools.  The program's build process requires [`goimports`](https://pkg.go.dev/golang.org/x/tools/cmd/goimports), [`stringer`](https://pkg.go.dev/golang.org/x/tools/cmd/stringer), and the [Pigeon parser generator](https://pkg.go.dev/github.com/PuerkitoBio/pigeon).  These can be installed from the command line with
```bash
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/tools/cmd/stringer@latest
go install github.com/mna/pigeon@latest
```
Once these dependencies are satisfied, `edif2qmasm` can be rebuilt with
```bash
go generate
go build
go install
```

or
```bash
make
make install
```
