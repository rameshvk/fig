package parse

import (
	"strconv"
	"unicode"
)

type tokenizer struct {
	seen  []rune
	start int
}

func (t *tokenizer) next(r rune, offset int, last bool, errs *[]error) (string, int, int, bool) {
	if len(t.seen) == 0 && unicode.IsSpace(r) {
		return "", -1, -1, last
	}

	if len(t.seen) == 0 {
		t.start = offset
	}

	t.seen = append(t.seen, r)
	if t.seen[0] == '"' {
		return t.quote(offset, last, errs)
	}
	if unicode.IsDigit(t.seen[0]) {
		return t.number(offset, last, errs)
	}
	if _, ok := priority[string(t.seen[:1])]; ok {
		return t.operator(offset, last, errs)
	}
	return t.name(offset, last, errs)
}

func (t *tokenizer) quote(offset int, last bool, errs *[]error) (string, int, int, bool) {
	l := len(t.seen)
	if l == 1 || (l > 1 && t.seen[l-1] != '"') || (l > 2 && t.seen[l-2] == '\\') {
		if !last {
			return "", -1, -1, false
		}
		// incomplete
		*errs = append(*errs, IncompleteStringError(t.start))
		if t.seen[l-1] == '\\' {
			t.seen = append(t.seen, '\\')
		}
		t.seen = append(t.seen, '"')
	}

	var result string
	result, t.seen = string(t.seen), t.seen[:0]
	return result, t.start, t.start + len(result), true
}

func (t *tokenizer) number(offset int, last bool, errs *[]error) (string, int, int, bool) {
	_, err := strconv.ParseFloat(string(t.seen), 64)
	if err == nil {
		if last {
			result, start := string(t.seen), t.start
			t.seen = t.seen[:0]
			return result, start, start + len(result), true
		}
		return "", -1, -1, false
	}

	l := len(t.seen)
	result, start := string(t.seen[:l-1]), t.start
	if unicode.IsSpace(t.seen[l-1]) {
		t.seen = t.seen[:0]
	} else {
		t.seen = t.seen[l-1:]
		t.start = offset
	}
	return result, start, start + len(result), true
}

func (t *tokenizer) operator(offset int, last bool, errs *[]error) (string, int, int, bool) {
	_, ok := priority[string(t.seen)]
	if ok {
		if last {
			result, start := string(t.seen), t.start
			t.seen = t.seen[:0]
			return result, start, start + len(result), true
		}
		return "", -1, -1, false
	}

	l := len(t.seen)
	result, start := string(t.seen[:l-1]), t.start
	if unicode.IsSpace(t.seen[l-1]) {
		t.seen = t.seen[:0]
	} else {
		t.seen = t.seen[l-1:]
		t.start = offset
	}
	return result, start, start + len(result), true
}

func (t *tokenizer) name(offset int, last bool, errs *[]error) (string, int, int, bool) {
	l := len(t.seen)
	if l == 1 && !unicode.IsLetter(t.seen[0]) {
		*errs = append(*errs, InvalidCharacterError(offset))
		t.seen = t.seen[:0]
		return "", -1, -1, false
	}

	// TODO: tighten this up!
	if _, ok := priority[string(t.seen[l-1:])]; ok || unicode.IsSpace(t.seen[l-1]) {
		result, start := string(t.seen[:l-1]), t.start

		if !ok {
			t.seen = t.seen[:0]
		} else {
			t.seen = t.seen[l-1:]
			t.start = offset
		}
		return result, start, start + len(result), true
	}

	if last {
		result, start := string(t.seen), t.start
		t.seen = t.seen[:0]
		return result, start, start + len(result), true
	}

	return "", -1, -1, false
}
