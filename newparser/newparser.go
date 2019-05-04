package newparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// NodeType is a type of node
type NodeType int

// Type of node
const (
	TextNode NodeType = iota
	ElementNode
	RootNode
)

// Node represents a text node or element node of HTML
type Node struct {
	NodeType   NodeType
	Data       string
	Children   []*Node
	Attributes map[string]string
}

// NewTextNode creates it
func NewTextNode(s string) *Node {
	return &Node{TextNode, s, []*Node{}, make(map[string]string)}
}

// NewElementNode creates it
func NewElementNode(tagName string, attrs map[string]string, ch []*Node) *Node {
	return &Node{ElementNode, tagName, ch, attrs}
}

// NewRootNode creates a node without attributes and tag name
func NewRootNode(ch []*Node) *Node {
	return &Node{RootNode, "", ch, make(map[string]string)}
}

// TagName returns tag of element node, empty if not element node
func (n *Node) TagName() string {
	return n.Data
}

///////----------------
///////----------------
///////----------------
///////----------------
///////----------------
///////----------------
///////----------------
///////----------------
///////----------------
///////----------------
///////----------------
///////----------------
///////----------------
///////----------------
///////----------------
///////----------------

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

type parser struct {
	t *tokenizer
}

func newParser(r io.Reader) *parser {
	return &parser{newTokenizer(r)}
}

func (p *parser) parse() (*Node, error) {
	t, err := p.t.readToken()
	if t.tokType == text {
		return NewTextNode(t.s), err
	} else if t.tokType == startTag {
		n, err := p.parseElement(t)
		return n, err
	}
	return nil, nil
}

func (p *parser) parseElement(t token) (*Node, error) {
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
