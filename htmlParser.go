package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type tokenType int

const (
	eof tokenType = iota
	startTag
	endTag
	text
)

type token struct {
	tokType tokenType
	s       string
	attrs   map[string]string
}

type tokenizer struct {
	r *bufio.Reader
}

func newTokenizer(r io.Reader) *tokenizer {
	return &tokenizer{bufio.NewReader(r)}
}

var defTok = token{}

func (t *tokenizer) readRune() (c rune, err error) {
	c, _, err = t.r.ReadRune()
	return
}

func (t *tokenizer) readToken() (tok token, err error) {
	c, err := t.readRune()
	if err == io.EOF {
		tok.tokType = eof
		err = nil
		return
	} else if err != nil {
		return defTok, err
	}

	switch c {
	case '<':
		tok, err := t.readTag()
		return tok, err
	}

	t.r.UnreadRune()
	tok.s, err = t.readUntil('<')
	if err != nil && err != io.EOF {
		return
	}
	t.r.UnreadRune()
	tok.tokType = text
	err = nil
	return
}

func (t *tokenizer) readTag() (tok token, err error) {
	c, err := t.readRune()
	if err != nil {
		return
	}
	if c == '/' {
		tok, err = t.readEndTag()
		return
	}
	tok.s, err = t.readTagName()
	if err != nil {
		return
	}
	tok.attrs, err = t.readAttrs()
	tok.tokType = startTag
	return
}

func (t *tokenizer) readEndTag() (tok token, err error) {
	c, err := t.readRune()
	if err != nil {
		return defTok, err
	}
	n, err := t.readTagName()
	if err != nil {
		return defTok, err
	}
	c, err = t.readRune()
	if c != '>' {
		return defTok, fmt.Errorf("CLOSE END TAG")
	}
	tok.tokType = endTag
	tok.s = n
	return
}

func (t *tokenizer) readTagName() (string, error) {
	s := new(strings.Builder)
	t.r.UnreadByte()
	for {
		c, err := t.readRune()
		if err != nil {
			return "", err
		}
		if c == '>' || unicode.IsSpace(c) || c == '/' {
			t.r.UnreadRune()
			return s.String(), nil
		}
		c = unicode.ToLower(c)
		s.WriteRune(c)
	}
}

func (t *tokenizer) readAttrs() (attrs map[string]string, err error) {
	attrs = make(map[string]string)
	var c rune
	var k, v string
	for {
		c, err = t.readRune()
		if err != nil {
			return
		}
		if unicode.IsSpace(c) {
			continue
		}
		if c == '>' {
			return
		}
		if c == '/' {
			err = t.consumeChar('>')
			return
		}
		t.r.UnreadRune()
		k, err = t.readUntil('=')
		if err != nil {
			return
		}
		err = t.consumeChar('"')
		if err != nil {
			return
		}
		v, err = t.readUntil('"')
		if err != nil {
			return
		}
		attrs[k] = v
	}
}

func (t *tokenizer) consumeChar(c rune) error {
	c2, err := t.readRune()
	if err != nil {
		return err
	} else if c2 != c {
		return fmt.Errorf("expected %c but got %c", c, c2)
	} else {
		return nil
	}
}

func (t *tokenizer) readUntil(fin rune) (string, error) {
	s := new(strings.Builder)
	for {
		c, err := t.readRune()
		if err != nil {
			return s.String(), err
		}
		if c == fin {
			return s.String(), nil
		}
		s.WriteRune(c)
	}
}

func (t token) String() string {
	return fmt.Sprintf("(%s, %q)", t.tokType, t.s)
}

type htmlParser struct {
	t *tokenizer
}

func newParser(r io.Reader) *htmlParser {
	return &htmlParser{newTokenizer(r)}
}

// parseHTML parses restricted subset of HTML
func parseHTML(r io.Reader) (*Node, error) {
	n, err := newParser(r).parse()
	return n, err
}

type concatReader struct {
	r1     io.Reader
	r2     io.Reader
	first  bool
	second bool
}

func (r *concatReader) Read(p []byte) (int, error) {
	if !r.second {
		return 0, io.EOF
	}
	if !r.first {
		n, err := r.r2.Read(p)
		if err == io.EOF {
			r.second = false
			return n, nil
		} else if err != nil {
			return n, err
		}
		return n, err
	}
	n, err := r.r1.Read(p)
	if err == io.EOF {
		r.first = false
		return n, nil
	} else if err != nil {
		return n, err
	}
	return n, err
}

func newConcatReader(r1, r2 io.Reader) io.Reader {
	return &concatReader{r1, r2, true, true}
}

// parseHTMLWrapped makes adds imageneary root element,
// so that strings like "<p></p><div></div>" could be parsed
func parseHTMLWrapped(r io.Reader) (*Node, error) {
	nodes := []*Node{}
	p := newParser(r)
	var err error
	for {
		n, err := p.parse()
		if err == nil && n != nil {
			nodes = append(nodes, n)
		} else {
			break
		}
	}

	return NewRootNode(nodes), err
}

func (p *htmlParser) parse() (*Node, error) {
	t, err := p.t.readToken()
	if t.tokType == text {
		return NewTextNode(t.s), err
	} else if t.tokType == startTag {
		n, err := p.parseElement(t)
		return n, err
	}
	return nil, nil
}

func (p *htmlParser) parseElement(t token) (*Node, error) {
	tagName := t.s
	attrs := t.attrs
	children := []*Node{}
	for {
		t, err := p.t.readToken()
		if err != nil {
			return nil, err
		}
		if t.tokType == text {
			children = append(children, NewTextNode(t.s))
		} else if t.tokType == startTag {
			n, err := p.parseElement(t)
			if err != nil {
				return nil, err
			}
			children = append(children, n)
		} else if t.tokType == endTag {
			if tagName != t.s {
				return nil, fmt.Errorf("Wrong end tag name")
			}
			return NewElementNode(tagName, attrs, children), nil
		}
	}
}
