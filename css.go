package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Example:
// h1, h2, h3 { margin: auto; color: #cc0000; }
// div.note { margin-bottom: 20px; padding: 10px; }
// #answer { display: none; }

type Stylesheet struct {
	rules []Rule
}

type Rule struct {
	selector   []Selector
	declarator []Declarator
}

type Selector struct {
	tagName *string
	id      *string
	class   []string
}

type Declarator struct {
	name      string
	valueType ValueType
	value     Value
}

type ValueType int

const (
	Keyword ValueType = iota
	Length
	ColorValue
)

type UnitType int

const (
	Px UnitType = iota
)

type Value struct {
	keyword  string
	length   float32
	unitType UnitType
	color    Colorr
}

type Colorr struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

// ParseStylesheet should parse CSS stylesheet
func ParseStylesheet(r io.Reader) (*Stylesheet, error) {
	reader := bufio.NewReader(r)
	rules := []Rule{}
	for {
		rule, err := parseRule(reader)
		if err != nil {
			return nil, err
		} else if rule == nil {
			break
		} else {
			rules = append(rules, *rule)
		}
	}
	return &Stylesheet{rules}, nil
}

func parseRule(r *bufio.Reader) (*Rule, error) {
	selectors, err := parseSelectors(r)
	if err != nil {
		return nil, err
	}
	if selectors == nil || len(selectors) == 0 {
		return nil, nil
	}
	declarators, err := parseDeclarators(r)
	if err != nil {
		return nil, err
	}
	return &Rule{selectors, declarators}, nil
}

func parseSelectors(r *bufio.Reader) ([]Selector, error) {
	selectors := []Selector{}
	for {
		sel, err := parseSelector(r)
		if sel == nil {
			return selectors, nil
		} else if err != nil {
			return nil, err
		}
		selectors = append(selectors, *sel)
		skipSpaces(r)
		if !isNextChar(r, ',') {
			break
		}
		r.ReadRune()
	}
	return selectors, nil
}

var nameStart = regexp.MustCompile("[_a-z]")
var nameChar = regexp.MustCompile("[_a-z0-9-]")

// parseSelector returns selector,
// or nil as selector if parse selector is not possible
// or returns an error
func parseSelector(r *bufio.Reader) (*Selector, error) {
	var selector *Selector
	skipSpaces(r)
	for {
		c, _, err := r.ReadRune()
		if err == io.EOF {
			return selector, nil
		} else if err != nil {
			return selector, err
		} else if c == '*' {
			// universal
		} else if c == '#' {
			name, err := readName(r)
			if err != nil {
				return selector, err
			}
			if selector == nil {
				selector = &Selector{}
			}
			selector.id = &name
		} else if c == '.' {
			name, err := readName(r)
			if err != nil {
				return selector, err
			}
			if selector == nil {
				selector = &Selector{}
			}
			selector.class = append(selector.class, name)
		} else if nameStart.MatchString(string(c)) {
			r.UnreadRune()
			name, err := readName(r)
			if err != nil && err != io.EOF {
				return selector, err
			}
			if selector == nil {
				selector = &Selector{}
			}
			selector.tagName = &name
		} else {
			r.UnreadRune()
			return selector, nil
		}
	}
}

func parseDeclarators(r *bufio.Reader) ([]Declarator, error) {
	skipSpaces(r)
	declarators := []Declarator{}
	if !consumeRequired(r, '{') {
		return declarators, fmt.Errorf("{ is required")
	}
	for {
		if isNextChar(r, ';') {
			r.ReadRune()
		}
		d, err := parseDeclarator(r)
		if err != nil {
			return nil, err
		} else if d == nil {
			break
		} else {
			declarators = append(declarators, *d)
		}
	}
	if !consumeRequired(r, '}') {
		return declarators, fmt.Errorf("} is required")
	}
	return declarators, nil
}

func parseDeclarator(r *bufio.Reader) (*Declarator, error) {
	skipSpaces(r)
	if !isNextCharMatches(r, nameStart) {
		return nil, nil
	}
	name, err := readName(r)
	if err != nil {
		return nil, err
	}
	declarator := new(Declarator)
	declarator.name = name
	if !isNextChar(r, ':') {
		return nil, fmt.Errorf("NEXT CHAR SHOULD BE COLON :")
	}
	r.ReadRune()
	skipSpaces(r)
	if isNextCharMatches(r, nameStart) {
		// keyword
		keyword, err := readName(r)
		if err != nil {
			return nil, err
		}
		declarator.valueType = Keyword
		declarator.value.keyword = keyword
	} else if isNextChar(r, '#') {
		// color
		color, err := readColor(r)
		if err != nil {
			return nil, err
		}
		declarator.valueType = ColorValue
		declarator.value.color = color
	} else {
		// length
		length, err := readLength(r)
		if err != nil {
			return nil, err
		}
		declarator.valueType = Length
		declarator.value.unitType = Px
		declarator.value.length = float32(length)
	}
	return declarator, nil
}

func readLength(r *bufio.Reader) (float64, error) {
	s := ""
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return 0.0, err
		} else if unicode.IsNumber(c) {
			s += string(c)
		} else {
			r.UnreadRune()
			break
		}
	}
	matches, err := matchStringInsens(r, "px")
	if err != nil {
		return 0.0, err
	} else if !matches {
		return 0.0, fmt.Errorf("px is expected")
	}
	return strconv.ParseFloat(s, 32)
}

func matchStringInsens(r *bufio.Reader, s string) (bool, error) {
	for _, c := range s {
		x, _, err := r.ReadRune()
		if err != nil {
			return false, err
		} else if unicode.ToUpper(x) != unicode.ToUpper(c) {
			return false, nil
		}
	}
	return true, nil
}

var hexDigit = regexp.MustCompile("[0-9a-fA-F]")

// readColor parses color, must be # folowed by six hex digits
func readColor(r *bufio.Reader) (Colorr, error) {
	r.ReadRune()
	var col Colorr
	s := ""
	for {
		if isNextCharMatches(r, hexDigit) {
			s += string(getChar(r))
		} else {
			break
		}
	}
	if len(s) != 6 {
		return col, fmt.Errorf("color must have six digits")
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return col, err
	}
	col.r = b[0]
	col.g = b[1]
	col.b = b[2]
	return col, nil
}

func getChar(r *bufio.Reader) rune {
	c, _, _ := r.ReadRune()
	return c
}

func readName(r *bufio.Reader) (string, error) {
	s := new(strings.Builder)
	c, _, err := r.ReadRune()
	if err != nil {
		return "", err
	}
	if !nameStart.MatchString(string(c)) {
		return "", fmt.Errorf("expected name start character bug got %q", c)
	}
	s.WriteRune(c)
	for {
		c, _, err := r.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", nil
		} else if nameChar.MatchString(string(c)) {
			s.WriteRune(c)
		} else {
			r.UnreadRune()
			break
		}
	}
	return s.String(), nil
}

func skipSpaces(r *bufio.Reader) error {
	for {
		b, _, err := r.ReadRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(b) {
			r.UnreadRune()
			return nil
		}
	}
}

func isNextChar(r *bufio.Reader, c rune) bool {
	x, _, err := r.ReadRune()
	r.UnreadRune()
	return err == nil && x == c
}

func isNextCharMatches(r *bufio.Reader, regex *regexp.Regexp) bool {
	x, _, err := r.ReadRune()
	r.UnreadRune()
	return err == nil && regex.MatchString(string(x))
}

func peekAndUnread(r *bufio.Reader) rune {
	x, _, _ := r.ReadRune()
	r.UnreadRune()
	return x
}

func consumeRequired(r *bufio.Reader, c rune) bool {
	x, _, err := r.ReadRune()
	return err == nil && c == x
}
