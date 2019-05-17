package main

import (
	"strings"
	"testing"
)

func TestNodesToBoxes(t *testing.T) {
	type args struct {
		node *styledNode
	}

	css := `p {display: block} em {display: inline} strong {display: inline}`
	cssR := strings.NewReader(css)

	t.Run("", func(t *testing.T) {
		sn := makeStyledNodeFromString(strings.NewReader("1"), cssR)
		box := nodesToBoxes(sn)
		if box.boxType != rootBox {
			t.Errorf("should be root box")
		}
	})

	t.Run("", func(t *testing.T) {
		html := "<p>1</p>"
		sn := makeStyledNodeFromString(strings.NewReader(html), cssR)
		box := nodesToBoxes(sn)
		if !(len(box.children) == 1 && box.children[0].tag() == "p") {
			t.Errorf("should have one <p>")
		}
	})

	t.Run("", func(t *testing.T) {
		html := "<p>1</p><p>2</p>"
		sn := makeStyledNodeFromString(strings.NewReader(html), cssR)
		box := nodesToBoxes(sn)
		if !(len(box.children) == 2 && box.children[0].tag() == "p" && box.children[1].tag() == "p") {
			t.Errorf("should have two <p>")
		}
	})

	t.Run("", func(t *testing.T) {
		// 		var html = `<P>Several <EM>emphasized words</EM> appear
		// <STRONG>in this</STRONG> sentence, dear.</P>`
		// 		sn := makeStyledNodeFromString(strings.NewReader(html), cssR)
		// 		box := nodesToBoxes(sn)
		// 		if !(len(box.children) == 1 && box.children[0].tag() == "p") {
		// 			t.Errorf("should have one <p>")
		// 		}
	})

	t.Run("", func(t *testing.T) {
		var html = "<p>1</p><span>2</span>"
		sn := makeStyledNodeFromString(strings.NewReader(html), cssR)
		box := nodesToBoxes(sn)
		if !(len(box.children) == 2 && box.children[0].tag() == "p" && box.children[1].isAnonymous()) {
			t.Errorf("should have one <p>")
		}
	})

	t.Run("", func(t *testing.T) {
		var html = "<span>2</span><p>1</p>"
		sn := makeStyledNodeFromString(strings.NewReader(html), cssR)
		box := nodesToBoxes(sn)
		if !(len(box.children) == 2 && box.children[1].tag() == "p" && box.children[0].isAnonymous()) {
			t.Errorf("should have one <p>")
		}
	})

	t.Run("", func(t *testing.T) {
		var html = "<span>3</span><span>2</span><p>1</p>"
		sn := makeStyledNodeFromString(strings.NewReader(html), cssR)
		box := nodesToBoxes(sn)
		if !(len(box.children) == 2 && box.children[1].tag() == "p" && box.children[0].isAnonymous()) {
			t.Errorf("should have one <p>")
		}
	})

	t.Run("", func(t *testing.T) {
		var html = "<span>3</span><span>2</span>"
		sn := makeStyledNodeFromString(strings.NewReader(html), cssR)
		box := nodesToBoxes(sn)
		if !(len(box.children) == 2 && box.children[1].tag() == "span" && !box.children[0].isAnonymous()) {
			t.Errorf("should have one <p>")
		}
	})

	t.Run("", func(t *testing.T) {
		var html = "<span>3</span><p>1</p><span>2</span>"
		sn := makeStyledNodeFromString(strings.NewReader(html), cssR)
		box := nodesToBoxes(sn)
		if !(len(box.children) == 3 && box.children[0].isAnonymous() && box.children[1].tag() == "p" && box.children[2].isAnonymous() == true) {
			t.Errorf("should be 3 blocks")
		}
	})
}
