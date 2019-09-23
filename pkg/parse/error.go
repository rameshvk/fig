package parse

import "strconv"

// ParseError represents a parse error
type ParseError interface {
	ErrorOffset() int
}

// MismatchedBracesError is when open and close braces/parens
// don't match
type MismatchedBracesError int

func (e MismatchedBracesError) Error() string {
	return "mismatched braces/parens at " + strconv.Itoa(int(e))
}
func (e MismatchedBracesError) ErrorOffset() int {
	return int(e)
}
func (e MismatchedBracesError) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.Error() + `"`), nil
}

// IncompleteBracesError is when terminating braces/parens
// are missing
type IncompleteBracesError int

func (e IncompleteBracesError) Error() string {
	return "incomplete braces/parens at " + strconv.Itoa(int(e))
}
func (e IncompleteBracesError) ErrorOffset() int {
	return int(e)
}
func (e IncompleteBracesError) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.Error() + `"`), nil
}

// IncompleteStringError is when a string is not terminated
type IncompleteStringError int

func (e IncompleteStringError) Error() string {
	return "unterminated string at " + strconv.Itoa(int(e))
}
func (e IncompleteStringError) ErrorOffset() int {
	return int(e)
}
func (e IncompleteStringError) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.Error() + `"`), nil
}

// InvalidCharacterError is an unexpected character is found
type InvalidCharacterError int

func (e InvalidCharacterError) Error() string {
	return "invalid character at " + strconv.Itoa(int(e))
}
func (e InvalidCharacterError) ErrorOffset() int {
	return int(e)
}
func (e InvalidCharacterError) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.Error() + `"`), nil
}

// MissingOperatorError is when two terms are adjacent without an
// operator in between
type MissingOperatorError int

func (e MissingOperatorError) Error() string {
	return "missing operator at " + strconv.Itoa(int(e))
}
func (e MissingOperatorError) ErrorOffset() int {
	return int(e)
}
func (e MissingOperatorError) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.Error() + `"`), nil
}

// MissingTermError is when a term is missing altogether
type MissingTermError int

func (e MissingTermError) Error() string {
	return "missing term at " + strconv.Itoa(int(e))
}
func (e MissingTermError) ErrorOffset() int {
	return int(e)
}
func (e MissingTermError) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.Error() + `"`), nil
}
