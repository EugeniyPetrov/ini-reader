package ini_reader

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var ErrUnexpectedEOF = errors.New("unexpected EOF")

type ParseError struct {
	line int
	err  error
}

func NewParseError(line int, err error) *ParseError {
	return &ParseError{
		line: line,
		err:  err,
	}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("ini: line:%d %s", e.line, e.err)
}

func (e *ParseError) Unwrap() error {
	return e.err
}

func (e *ParseError) Line() int {
	return e.line
}

type stateFn func() stateFn

var eof rune = -1

type sectionIterator func(*Section, error) bool

type Section struct {
	Name       string
	Properties map[string]any
}

type Reader struct {
	line    int
	atEOF   bool
	input   *bufio.Reader
	section *Section
	err     error
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		line:  1,
		input: bufio.NewReader(r),
	}
}

func (r *Reader) errorf(format string, args ...any) {
	format = fmt.Sprintf("ini: line:%d %s", r.line, format)
	panic(
		NewParseError(r.line, fmt.Errorf(format, args...)),
	)
}

func (r *Reader) error(err error) {
	r.errorf("%w", err)
}

func (r *Reader) next() rune {
	rn, _, err := r.input.ReadRune()
	if err == io.EOF {
		r.atEOF = true
		return eof
	}
	if err != nil {
		r.errorf("error reading rune. %w", err)
	}
	if rn == '\n' {
		r.line++
	}
	return rn
}

func (r *Reader) backup() {
	if r.atEOF {
		return
	}

	err := r.input.UnreadRune()
	if err != nil {
		r.errorf("error unread rune. %w", err)
	}

	// @todo decrease line number if unread rune is '\n'
}

func (r *Reader) nextLine() {
	for {
		rn := r.next()
		if rn == eof || rn == '\n' {
			return
		}
	}
}

func (r *Reader) accept(valid string) bool {
	if strings.ContainsRune(valid, r.next()) {
		return true
	}
	r.backup()
	return false
}

func (r *Reader) parsePropertyStringValue(propertyName string, quote rune) stateFn {
	return func() stateFn {
		value := bytes.NewBufferString("")
		for {
			sym := r.next()

			switch {
			case sym == eof:
				r.error(ErrUnexpectedEOF)
			case sym == quote:
				r.section.Properties[propertyName] = value.String()
				r.accept("\n")
				return r.parseProperties
			default:
				value.WriteRune(sym)
			}
		}
	}
}

func (r *Reader) isInteger(value string) bool {
	if value == "" {
		return false
	}

	if len(value) > 1 && value[0] == '0' {
		return false
	}

	for _, c := range value[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (r *Reader) isFloat(value string) bool {
	for _, c := range value {
		if c < '0' || c > '9' || c != '.' {
			return false
		}
	}
	return true
}

func (r *Reader) scalarValue(value string) any {
	if value == "" {
		return nil
	} else if value == "true" {
		return true
	} else if value == "false" {
		return false
	} else if r.isInteger(value) {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			r.error(fmt.Errorf("error parsing integer value %q. %w", value, err))
		}
		return v
	} else if r.isFloat(value) {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			r.error(err)
		}
		return v
	}

	return value
}

func (r *Reader) parsePropertyScalarValue(propertyName string) stateFn {
	return func() stateFn {
		value := bytes.NewBufferString("")
		for {
			sym := r.next()

			switch {
			case sym == '\n' || sym == eof || sym == ';':
				r.section.Properties[propertyName] = r.scalarValue(
					strings.TrimSpace(value.String()),
				)
				if sym == ';' {
					r.nextLine()
				}
				return r.parseProperties
			default:
				value.WriteRune(sym)
			}
		}
	}
}

func (r *Reader) parsePropertyValue(propertyName string) stateFn {
	r.section.Properties[propertyName] = nil

	return func() stateFn {
		for {
			sym := r.next()

			switch {
			case sym == eof:
				return nil
			case sym == ' ':
				continue
			case sym == '"' || sym == '\'':
				return r.parsePropertyStringValue(propertyName, sym)
			default:
				r.backup()
				return r.parsePropertyScalarValue(propertyName)
			}
		}
	}
}

func (r *Reader) isWhitespace(sym rune) bool {
	return sym == ' ' || sym == '\t'
}

func (r *Reader) isPropertyName(sym rune) bool {
	return sym != '=' && sym != ';'
}

func (r *Reader) parsePropertyName() stateFn {
	name := bytes.NewBufferString("")
	for {
		sym := r.next()

		switch {
		case sym == eof:
			return nil
		case r.isPropertyName(sym):
			name.WriteRune(sym)
		case sym == '=':
			return r.parsePropertyValue(strings.TrimSpace(name.String()))
		default:
			r.errorf("unexpected symbol %q", sym)
		}
	}
}

func (r *Reader) parseProperties() stateFn {
	for {
		sym := r.next()

		switch {
		case sym == eof:
			return nil
		case sym == '\n':
			return nil
		case r.isWhitespace(sym):
			continue
		case sym == ';':
			r.nextLine()
			continue
		case r.isPropertyName(sym):
			r.backup()
			return r.parsePropertyName
		default:
			r.errorf("unexpected symbol %q", sym)
		}
	}
}

func (r *Reader) newSection() {
	r.section = &Section{
		Properties: make(map[string]any),
	}
}

func (r *Reader) parseSectionHeader() stateFn {
	r.newSection()

	name := bytes.NewBufferString("")
	for {
		sym := r.next()

		switch sym {
		case eof:
			return nil
		case ']':
			r.section.Name = name.String()
			r.accept("\n")
			return r.parseProperties
		default:
			name.WriteRune(sym)
		}
	}
}

func (r *Reader) nextSection() stateFn {
	for {
		sym := r.next()
		switch {
		case sym == eof:
			return nil
		case r.isWhitespace(sym):
			continue
		case sym == ';':
			r.nextLine()
		case sym == '\n':
			continue
		case sym == '[':
			return r.parseSectionHeader
		case r.isPropertyName(sym):
			r.backup()
			r.newSection()
			return r.parseProperties
		}
	}
}

func (r *Reader) Next() bool {
	defer func() {
		e := recover()
		if e != nil {
			r.err = e.(error)
		}
	}()

	if r.atEOF {
		return false
	}

	state := r.nextSection
	for {
		state = state()
		if state == nil {
			return !(r.atEOF && r.section == nil)
		}
	}
}

func (r *Reader) Section() *Section {
	return r.section
}

func (r *Reader) Err() error {
	return r.err
}

func (r *Reader) ReadAll() ([]*Section, error) {
	sections := make([]*Section, 0)
	for r.Next() {
		sections = append(sections, r.Section())
	}
	err := r.Err()
	if err != nil {
		return nil, err
	}

	return sections, nil
}
