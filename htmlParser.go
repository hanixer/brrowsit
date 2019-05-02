package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
)

// The following syntax is allowed:

//     Balanced tags: <p>...</p>
//     Attributes with quoted values: id="main"
//     Text nodes: <em>world</em>

type TokenType int

const (
	less TokenType = iota
	greater
	lessSlash
	text
	ws
	equal
	stringTok
	eof
)

type Token struct {
	Type TokenType
	Data string
}

type Tokenizer struct {
	r      *bufio.Reader
	cached *Token
}

func (t Token) String() string {
	return fmt.Sprintf("(%s, %q)", t.Type, t.Data)
}

func (s *Tokenizer) IsNextType(typ TokenType) bool {
	if s.cached == nil {
		t := s.Scan()
		s.cached = &t
	}
	return s.cached.Type == typ
}

func (s *Tokenizer) Scan() Token {
	var token Token
	if s.cached != nil {
		token = *s.cached
		s.cached = nil
		return token
	}
	for {
		c, _, err := s.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				token.Type = eof
				break
			} else {
				log.Fatalln("Scan:", err)
			}
		}

		if c == '<' {
			c, _, err := s.r.ReadRune()
			if err != nil {
				log.Fatalln("Scan 2:", err)
			}
			if c == '/' {
				token.Type = lessSlash
				break
			}
			s.r.UnreadRune()
			token.Type = less
			break
		} else if c == '>' {
			token.Type = greater
			break
		} else if c == '=' {
			token.Type = equal
			break
		} else if c == '"' {
			token.Type = stringTok
			token.Data = scanStringToken(s.r)
			break
		} else if unicode.IsSpace(rune(c)) {
			token.Type = ws
			s.r.UnreadRune()
			token.Data = scanWhitespace(s.r)
			break
		} else {
			token.Type = text
			s.r.UnreadRune()
			token.Data = scanText(s.r)
			break
		}
	}
	return token
}

func (t *Tokenizer) Unscan(tok Token) {
	t.cached = &tok
}

func scanStringToken(r *bufio.Reader) string {
	builder := new(strings.Builder)
	for {
		b, _, err := r.ReadRune()
		if err != nil {
			log.Fatalln(err)
		}
		if b == '"' {
			break
		}
		builder.WriteRune(b)
	}
	return builder.String()
}

func scanWhitespace(r *bufio.Reader) string {
	builder := new(strings.Builder)
	for {
		b, _, err := r.ReadRune()
		if err != nil {
			log.Fatalln(err)
		}
		if !unicode.IsSpace(b) {
			r.UnreadRune()
			break
		}
		builder.WriteRune(b)
	}
	return builder.String()
}

func scanText(r *bufio.Reader) string {
	builder := new(strings.Builder)
	for {
		b, _, err := r.ReadRune()
		if err != nil {
			log.Fatalln(err)
		}
		if b == '<' || b == '>' || b == '/' || b == '=' || b == '"' || unicode.IsSpace(b) {
			r.UnreadRune()
			break
		}
		builder.WriteRune(b)
	}
	return builder.String()
}

// NewParser creates a new Tokenizer from Reader
func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{bufio.NewReader(r), nil}
}

type Parser struct {
	t *Tokenizer
}

// NewParser creates a new Parser from Reader
func NewParser(r io.Reader) *Parser {
	return &Parser{NewTokenizer(r)}
}

func ParseHtml(r io.Reader) *Node {
	p := NewParser(r)
	return p.Parse()
}

// Parse returns a single node or a single node if </ is encountered.
// Panic on other errors.
func (p *Parser) Parse() *Node {
	for {
		t := p.consumeSpaces()
		if t.Type == lessSlash {
			p.t.Unscan(t)
			return nil
		} else if t.Type == less {
			return p.parseElementNode()
		} else if t.Type == text {
			return NewTextNode(t.Data)
		} else {
			log.Fatalln("expected < or text, got", t)
		}
	}
}

func (p *Parser) parseElementNode() *Node {
	t := p.consumeSpaces()
	if t.Type != text {
		log.Fatalln("expected start tag name, got", t)
	}
	name := t.Data
	attrs := p.parseAttributes()
	childs := []*Node{}
	for {
		child := p.Parse()
		if child == nil {
			break
		}
		childs = append(childs, child)
	}
	t = p.consumeSpaces()
	t3 := p.consumeSpaces()
	t4 := p.consumeSpaces()
	if t.Type != lessSlash || t3.Type != text || t4.Type != greater {
		log.Fatalln("expected closing tag, got: ", t, t3, t4)
	}
	if t3.Data != name {
		log.Fatalln("unmatching closing tag name")
	}

	return NewElementNode(name, attrs, childs)
}

func (p *Parser) consumeSpaces() Token {
	for {
		t := p.t.Scan()
		if t.Type != ws {
			return t
		}
	}
}

func (p *Parser) parseAttributes() map[string]string {
	attrs := make(map[string]string)
	for {
		t := p.consumeSpaces()
		if t.Type == greater {
			break
		}
		if t.Type != text {
			log.Fatalln("expected key of attribute, got", t)
		}
		k := t.Data
		if p.consumeSpaces().Type != equal {
			log.Fatalln("expected equal sign, got")
		}
		t = p.consumeSpaces()
		if t.Type != stringTok {
			log.Fatalln("expected value of attribute, got", t)
		}
		v := t.Data
		attrs[k] = v
	}
	return attrs
}

func PrintNode(node *Node, w io.Writer) {
	printNode(node, w, 0)
}

func printNode(node *Node, w io.Writer, nesting int) {
	switch node.NodeType {
	case TextNode:
		io.WriteString(w, node.Data)
	case ElementNode:
		fmt.Fprintf(w, "<%s", node.Data)
		printAttributes(w, node.Attributes)
		fmt.Fprint(w, ">")
		newLinesNeccessary := len(node.Children) > 0 && node.Children[0].NodeType != TextNode
		for _, child := range node.Children {
			if newLinesNeccessary {
				printNesting(w, nesting+1)
			}
			printNode(child, w, nesting+1)
		}
		if newLinesNeccessary {
			printNesting(w, nesting)
		}
		fmt.Fprintf(w, "</%s>", node.Data)
	}
}

func printAttributes(w io.Writer, attrs map[string]string) {
	for k, v := range attrs {
		fmt.Fprintf(w, " %s=%q", k, v)
	}
}

func printNesting(w io.Writer, nesting int) {
	fmt.Fprintf(w, "\r\n")
	for count := 0; count < nesting; count++ {
		fmt.Fprintf(w, "  ")
	}
}

func (n *Node) String() string {
	builder := new(strings.Builder)
	PrintNode(n, builder)
	return builder.String()
}
