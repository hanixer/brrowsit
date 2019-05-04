package main

import (
	"strings"
	"testing"
)

func Test_tokenizer_readToken(t *testing.T) {
	tests := []struct {
		name string
		tzer *tokenizer
		want token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
	t.Run("empty - eof", func(t *testing.T) {
		tz := newTokenizer(strings.NewReader(""))
		tok, _ := tz.readToken()
		if tok.tokType != eof {
			t.Error("eof expected")
		}
	})
	t.Run("open tag", func(t *testing.T) {
		tz := newTokenizer(strings.NewReader("<TAG>"))
		tok, err := tz.readToken()
		expect := err == nil && tok.tokType == startTag && tok.s == "tag"
		if !expect {
			t.Error("error, or not start, or not lower: ", err, tok)
		}
	})
	t.Run("end tag", func(t *testing.T) {
		tz := newTokenizer(strings.NewReader("</TAG>"))
		tok, err := tz.readToken()
		expect := err == nil && tok.tokType == endTag && tok.s == "tag"
		if !expect {
			t.Error("error, or not end, or not lower: ", err, tok)
		}
	})
	t.Run("simple-attribute", func(t *testing.T) {
		tz := newTokenizer(strings.NewReader("<TAG k=\"v\">"))
		tok, err := tz.readToken()
		expect := err == nil && tok.tokType == startTag && tok.s == "tag" && tok.attrs != nil && tok.attrs["k"] == "v"
		if !expect {
			t.Error("error : ", err, tok)
		}
	})
	someText := "I am a simple text"
	t.Run("standalone-text", func(t *testing.T) {
		tz := newTokenizer(strings.NewReader(someText))
		tok, err := tz.readToken()
		expect := err == nil && tok.tokType == text && tok.s == someText
		if !expect {
			t.Error("error : ", err, tok)
		}
	})
	t.Run("text-inside-div", func(t *testing.T) {
		tz := newTokenizer(strings.NewReader("<div>" + someText + "</div>"))
		tok1, err1 := tz.readToken()
		tok2, err2 := tz.readToken()
		tok3, err3 := tz.readToken()
		expect := err1 == nil && err2 == nil && err3 == nil && tok1.tokType == startTag && tok2.tokType == text && tok2.s == someText && tok3.tokType == endTag
		if !expect {
			t.Error("error : ", err1, err2, err3, tok1, tok2, tok3)
		}
	})
}

func Test_parser_parse(t *testing.T) {
	someText := "I am a simple text"
	t.Run("standalone-text", func(t *testing.T) {
		p := newParser(strings.NewReader(someText))
		n, err := p.parse()
		expect := err == nil && n != nil && n.NodeType == TextNode
		if !expect {
			t.Error("error : ", err, n)
		}
	})
	t.Run("standalone-div", func(t *testing.T) {
		p := newParser(strings.NewReader("<div></div>"))
		n, err := p.parse()
		expect := err == nil && n != nil && n.NodeType == ElementNode && n.TagName() == "div"
		if !expect {
			t.Error("error : ", err, n)
		}
	})
	t.Run("div-with-text", func(t *testing.T) {
		p := newParser(strings.NewReader("<div>" + someText + "</div>"))
		n, err := p.parse()
		expect := err == nil && n != nil && n.NodeType == ElementNode && n.TagName() == "div" && len(n.Children) == 1 && n.Children[0].NodeType == TextNode && n.Children[0].Data == someText
		if !expect {
			t.Error("error : ", err, n)
		}
	})
	t.Run("div-with-div", func(t *testing.T) {
		p := newParser(strings.NewReader("<div><div></div></div>"))
		n, err := p.parse()
		expect := err == nil && n != nil && n.NodeType == ElementNode && n.TagName() == "div" && len(n.Children) == 1 && n.Children[0].NodeType == ElementNode && n.Children[0].Data == "div"
		if !expect {
			t.Error("error : ", err, n)
		}
	})
}
