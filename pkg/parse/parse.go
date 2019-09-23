// Package parse implements a simple expression parser.
//
// It returns a JSON version of the s-expression
package parse

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var priority = map[string]int{
	"(":  0,
	")":  0,
	"{":  0,
	"}":  0,
	",":  1,
	"=":  2,
	"|":  3,
	"&":  4,
	"==": 5,
	"!=": 5,
	"<":  7,
	">":  7,
	"<=": 7,
	">=": 7,
	"+":  10,
	"-":  10,
	"*":  20,
	"/":  20,
	"!":  30,
	".":  40,
}

// String parses the string and returns the parsed
// expression as a number, string or S-expression
func String(s string) (interface{}, []error) {
	var errs []error
	p := parsers{&parser{}}
	// The space hack is needed because the parse code below
	// does not handle finishing the token up properly
	result := p.parse(s+" ", &errs)
	result = normalize(result, &errs)
	return result, errs
}

type parsers []*parser

func (p *parsers) parse(s string, errs *[]error) interface{} {
	tok := &tokenizer{}
	l := len(s)

	for offset, r := range s {
		s, start, end, ok := tok.next(r, offset, offset+utf8.RuneLen(r) == l, errs)
		if !ok {
			continue
		}

		_, ok = priority[s]

		switch {
		case s == "":
			break
		case s == "(", s == "{":
			*p = append(*p, &parser{nesting: s, nestingStart: start})
		case s == ")", s == "}":
			nesting, nestingStart := p.top().nesting, p.top().nestingStart
			pair := nesting + s
			if pair != "()" && pair != "{}" {
				*errs = append(*errs, MismatchedBracesError(start))
				return p.finish(nestingStart, errs)
			}
			term := p.pop().finish(nestingStart, errs)
			term = p.top().term(pair, nestingStart, end, term)
			p.top().handleTerm(term, nestingStart, end, errs)
		case ok:
			p.top().handleOp(s, start, end, errs)
		default:
			p.top().handleTerm(s, start, end, errs)
		}
	}

	return p.finish(len(s), errs)
}

func (p *parsers) top() *parser {
	if len(*p) == 0 {
		return nil
	}
	return (*p)[len(*p)-1]
}

func (p *parsers) pop() *parser {
	l := len(*p)
	if l == 0 {
		return nil
	}
	top := (*p)[l-1]
	*p = (*p)[:l-1]
	return top
}

func (p *parsers) finish(start int, errs *[]error) interface{} {
	result := p.pop().finish(start, errs)
	for len(*p) > 0 {
		p.top().handleTerm(result, start, start, errs)
		result = p.pop().finish(start, errs)
	}
	return result
}

type parser struct {
	lastWasTerm  bool
	ops          []string
	starts, ends []int
	terms        []interface{}
	nesting      string
	nestingStart int
}

func (p *parser) term(op string, start, end int, terms ...interface{}) interface{} {
	loc := strconv.Itoa(start) + ":" + strconv.Itoa(end)
	return append([]interface{}{op, loc}, terms...)
}

func (p *parser) handleOp(op string, start, end int, errs *[]error) {
	if !p.lastWasTerm {
		p.terms = append(p.terms, MissingTermError(start))
	}

	pri := priority[op]
	isRightAssociative := op == "="
	for l := len(p.ops) - 1; l >= 0 && priority[p.ops[l]] >= pri; l-- {
		if isRightAssociative && p.ops[l] == op {
			break
		}
		p.unwindOp(errs)
	}

	p.ops = append(p.ops, op)
	p.starts = append(p.starts, start)
	p.ends = append(p.ends, end)
	p.lastWasTerm = false
}

func (p *parser) handleTerm(term interface{}, start, end int, errs *[]error) {
	term = p.wrapTerm(term, start, end)
	if p.lastWasTerm {
		if l := len(p.ops) - 1; l >= 0 && p.ops[l] == "." {
			p.unwindOp(errs)
		}

		term = p.term("", start, start, p.popTerm(start, errs), term)
	}
	p.terms = append(p.terms, term)
	p.lastWasTerm = true
}

func (p *parser) unwindOp(errs *[]error) {
	l := len(p.ops) - 1
	right := p.popTerm(p.starts[l], errs)
	left := p.popTerm(p.starts[l], errs)
	term := p.term(p.ops[l], p.starts[l], p.ends[l], left, right)
	p.terms = append(p.terms, term)
	p.ops, p.starts, p.ends = p.ops[:l], p.starts[:l], p.ends[:l]
}

func (p *parser) finish(start int, errs *[]error) interface{} {
	for l := len(p.ops) - 1; l >= 0; l-- {
		p.unwindOp(errs)
	}

	return p.popTerm(start, errs)
}

func (p *parser) popTerm(start int, err *[]error) (result interface{}) {
	if l := len(p.terms); l > 0 {
		result, p.terms = p.terms[l-1], p.terms[:l-1]
		return result
	}
	return MissingTermError(start)
}

func (p *parser) wrapTerm(t interface{}, start, end int) interface{} {
	s, ok := t.(string)
	if !ok {
		return t
	}

	if s[0] == '"' {
		var rs = make([]rune, 0, utf8.RuneCount([]byte(s))-2)
		var skip = false
		for _, r := range s[1 : len(s)-1] {
			if r != '\\' || skip {
				rs = append(rs, r)
				skip = false
			} else {
				skip = true
			}
		}
		return p.term("string", start, end, string(rs))
	}

	if unicode.IsDigit(([]rune(s))[0]) {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return p.term("number", start, end, f)
	}

	if x := strings.ToLower(s); x == "true" || x == "false" {
		return p.term("bool", start, end, x == "true")
	}

	return p.term("name", start, end, s)
}
