// This file is part of edif2qmasm.  It implements a parser for the Electronic
// Design Interchange Format (EDIF).  It basically just reads an EDIF file into
// a hierarchical format for later processing.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "TopLevel",
			pos:  position{line: 11, col: 1, offset: 274},
			expr: &actionExpr{
				pos: position{line: 11, col: 13, offset: 286},
				run: (*parser).callonTopLevel1,
				expr: &seqExpr{
					pos: position{line: 11, col: 13, offset: 286},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 11, col: 13, offset: 286},
							name: "Skip",
						},
						&labeledExpr{
							pos:   position{line: 11, col: 18, offset: 291},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 11, col: 20, offset: 293},
								name: "SExp",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 11, col: 25, offset: 298},
							name: "Skip",
						},
						&ruleRefExpr{
							pos:  position{line: 11, col: 30, offset: 303},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Digit",
			pos:  position{line: 20, col: 1, offset: 507},
			expr: &charClassMatcher{
				pos:        position{line: 20, col: 10, offset: 516},
				val:        "[\\p{Nd}]",
				classes:    []*unicode.RangeTable{rangeTable("Nd")},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 23, col: 1, offset: 591},
			expr: &charClassMatcher{
				pos:        position{line: 23, col: 15, offset: 605},
				val:        "[\\p{Zs}\\n\\r\\t]",
				chars:      []rune{'\n', '\r', '\t'},
				classes:    []*unicode.RangeTable{rangeTable("Zs")},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Skip",
			pos:  position{line: 26, col: 1, offset: 700},
			expr: &zeroOrMoreExpr{
				pos: position{line: 26, col: 9, offset: 708},
				expr: &ruleRefExpr{
					pos:  position{line: 26, col: 9, offset: 708},
					name: "Whitespace",
				},
			},
		},
		{
			name: "Symbol",
			pos:  position{line: 30, col: 1, offset: 862},
			expr: &actionExpr{
				pos: position{line: 30, col: 11, offset: 872},
				run: (*parser).callonSymbol1,
				expr: &seqExpr{
					pos: position{line: 30, col: 11, offset: 872},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 30, col: 11, offset: 872},
							val:        "[\\p{Lu}\\p{Ll}_]",
							chars:      []rune{'_'},
							classes:    []*unicode.RangeTable{rangeTable("Lu"), rangeTable("Ll")},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 30, col: 27, offset: 888},
							expr: &charClassMatcher{
								pos:        position{line: 30, col: 27, offset: 888},
								val:        "[^()\\p{Zs}\\n\\r\\t]",
								chars:      []rune{'(', ')', '\n', '\r', '\t'},
								classes:    []*unicode.RangeTable{rangeTable("Zs")},
								ignoreCase: false,
								inverted:   true,
							},
						},
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 35, col: 1, offset: 982},
			expr: &actionExpr{
				pos: position{line: 35, col: 12, offset: 993},
				run: (*parser).callonInteger1,
				expr: &seqExpr{
					pos: position{line: 35, col: 12, offset: 993},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 35, col: 12, offset: 993},
							expr: &charClassMatcher{
								pos:        position{line: 35, col: 12, offset: 993},
								val:        "[-+]",
								chars:      []rune{'-', '+'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 35, col: 18, offset: 999},
							expr: &ruleRefExpr{
								pos:  position{line: 35, col: 18, offset: 999},
								name: "Digit",
							},
						},
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 45, col: 1, offset: 1255},
			expr: &choiceExpr{
				pos: position{line: 45, col: 11, offset: 1265},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 45, col: 11, offset: 1265},
						run: (*parser).callonString2,
						expr: &seqExpr{
							pos: position{line: 45, col: 11, offset: 1265},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 45, col: 11, offset: 1265},
									val:        "\"",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 45, col: 15, offset: 1269},
									val:        "\\",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 45, col: 20, offset: 1274},
									expr: &charClassMatcher{
										pos:        position{line: 45, col: 20, offset: 1274},
										val:        "[^\"]",
										chars:      []rune{'"'},
										ignoreCase: false,
										inverted:   true,
									},
								},
								&litMatcher{
									pos:        position{line: 45, col: 26, offset: 1280},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 48, col: 5, offset: 1395},
						run: (*parser).callonString9,
						expr: &seqExpr{
							pos: position{line: 48, col: 5, offset: 1395},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 48, col: 5, offset: 1395},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 48, col: 9, offset: 1399},
									expr: &choiceExpr{
										pos: position{line: 48, col: 10, offset: 1400},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 48, col: 10, offset: 1400},
												name: "EscapedChar",
											},
											&charClassMatcher{
												pos:        position{line: 48, col: 24, offset: 1414},
												val:        "[^\"]",
												chars:      []rune{'"'},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 48, col: 31, offset: 1421},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 55, col: 1, offset: 1627},
			expr: &actionExpr{
				pos: position{line: 55, col: 16, offset: 1642},
				run: (*parser).callonEscapedChar1,
				expr: &seqExpr{
					pos: position{line: 55, col: 16, offset: 1642},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 55, col: 16, offset: 1642},
							val:        "\\",
							ignoreCase: false,
						},
						&anyMatcher{
							line: 55, col: 21, offset: 1647,
						},
					},
				},
			},
		},
		{
			name: "List",
			pos:  position{line: 71, col: 1, offset: 2068},
			expr: &actionExpr{
				pos: position{line: 71, col: 9, offset: 2076},
				run: (*parser).callonList1,
				expr: &seqExpr{
					pos: position{line: 71, col: 9, offset: 2076},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 71, col: 9, offset: 2076},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 71, col: 13, offset: 2080},
							name: "Skip",
						},
						&labeledExpr{
							pos:   position{line: 71, col: 18, offset: 2085},
							label: "ss",
							expr: &zeroOrOneExpr{
								pos: position{line: 71, col: 21, offset: 2088},
								expr: &seqExpr{
									pos: position{line: 71, col: 22, offset: 2089},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 71, col: 22, offset: 2089},
											name: "SExp",
										},
										&zeroOrMoreExpr{
											pos: position{line: 71, col: 27, offset: 2094},
											expr: &ruleRefExpr{
												pos:  position{line: 71, col: 27, offset: 2094},
												name: "AnotherSExp",
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 71, col: 42, offset: 2109},
							name: "Skip",
						},
						&litMatcher{
							pos:        position{line: 71, col: 47, offset: 2114},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SExp",
			pos:  position{line: 120, col: 1, offset: 3902},
			expr: &choiceExpr{
				pos: position{line: 120, col: 9, offset: 3910},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 120, col: 9, offset: 3910},
						name: "Symbol",
					},
					&ruleRefExpr{
						pos:  position{line: 120, col: 18, offset: 3919},
						name: "String",
					},
					&ruleRefExpr{
						pos:  position{line: 120, col: 27, offset: 3928},
						name: "Integer",
					},
					&ruleRefExpr{
						pos:  position{line: 120, col: 37, offset: 3938},
						name: "List",
					},
				},
			},
		},
		{
			name: "AnotherSExp",
			pos:  position{line: 123, col: 1, offset: 4007},
			expr: &actionExpr{
				pos: position{line: 123, col: 16, offset: 4022},
				run: (*parser).callonAnotherSExp1,
				expr: &seqExpr{
					pos: position{line: 123, col: 16, offset: 4022},
					exprs: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 123, col: 16, offset: 4022},
							expr: &ruleRefExpr{
								pos:  position{line: 123, col: 16, offset: 4022},
								name: "Whitespace",
							},
						},
						&labeledExpr{
							pos:   position{line: 123, col: 28, offset: 4034},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 123, col: 30, offset: 4036},
								name: "SExp",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 132, col: 1, offset: 4235},
			expr: &notExpr{
				pos: position{line: 132, col: 8, offset: 4242},
				expr: &anyMatcher{
					line: 132, col: 9, offset: 4243,
				},
			},
		},
	},
}

func (c *current) onTopLevel1(s interface{}) (interface{}, error) {
	sexp, ok := s.(EdifSExp)
	if !ok {
		return nil, fmt.Errorf("Failed to parse %q", c.text)
	}
	return sexp, nil
}

func (p *parser) callonTopLevel1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTopLevel1(stack["s"])
}

func (c *current) onSymbol1() (interface{}, error) {
	return EdifSymbol(c.text), nil
}

func (p *parser) callonSymbol1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSymbol1()
}

func (c *current) onInteger1() (interface{}, error) {
	num, err := strconv.Atoi(string(c.text))
	if err != nil {
		return nil, err
	}
	return EdifInteger(num), nil
}

func (p *parser) callonInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger1()
}

func (c *current) onString2() (interface{}, error) {
	// Verilog symbol beginning with a "\"
	return EdifString(c.text[1 : len(c.text)-1]), nil
}

func (p *parser) callonString2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString2()
}

func (c *current) onString9() (interface{}, error) {
	// String with ordinary character escapes
	return EdifString(c.text[1 : len(c.text)-1]), nil
}

func (p *parser) callonString9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString9()
}

func (c *current) onEscapedChar1() (interface{}, error) {
	switch c.text[1] {
	case '\\', '"':
		return c.text[1], nil
	case 'n':
		return '\n', nil
	case 't':
		return '\t', nil
	case 'r':
		return '\r', nil
	default:
		return "", fmt.Errorf("Unrecognized escape sequence \"\\%c\"", c.text[1])
	}
}

func (p *parser) callonEscapedChar1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEscapedChar1()
}

func (c *current) onList1(ss interface{}) (interface{}, error) {
	// Handle the trivial case first.
	if ss == nil {
		// Zero-element list
		return make(EdifList, 0), nil
	}

	// On a failed assertion in the remaining cases, return an internal
	// error to our parent.
	return func() (data interface{}, err error) {
		// Set up error handling.
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("Internal error parsing %q", c.text)
			}
		}()
		var ok bool
		checkAssert := func() {
			if !ok {
				panic("Internal edif.go parse error")
			}
		}

		// We expect to have either a one- or two-element list of
		// s-expressions.
		ssList, ok := ss.([]interface{})
		checkAssert()
		if len(ssList) == 0 {
			return nil, fmt.Errorf("Internal error: unexpected zero-element list")
		}
		ssHead, ok := ssList[0].(EdifSExp)
		checkAssert()
		if len(ssList) == 1 {
			// One-element list
			return []EdifSExp{ssHead}, nil
		}
		ssList, ok = ssList[1].([]interface{})
		checkAssert()
		sexps := make(EdifList, len(ssList)+1)
		sexps[0] = ssHead
		for i, s := range ssList {
			sexps[i+1], ok = s.(EdifSExp)
			checkAssert()
		}
		return sexps, nil
	}()
}

func (p *parser) callonList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onList1(stack["ss"])
}

func (c *current) onAnotherSExp1(s interface{}) (interface{}, error) {
	sexp, ok := s.(EdifSExp)
	if !ok {
		return nil, fmt.Errorf("Internal error parsing %#v", s)
	}
	return sexp, nil
}

func (p *parser) callonAnotherSExp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnotherSExp1(stack["s"])
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
