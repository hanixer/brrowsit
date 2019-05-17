package main

import (
	"strings"
	"testing"
)

func TestNodesToBoxes(t *testing.T) {
	type args struct {
		node *styledNode
	}

	t.Run("", func(t *testing.T) {
		sn := makeStyledNodeFromString(strings.NewReader("1"), strings.NewReader(""))
		box := nodesToBoxes(sn)
		if box.boxType != rootBox {
			t.Errorf("should be root box")
		}
	})

	t.Run("", func(t *testing.T) {
		html := "<p>1</p>"
		sn := makeStyledNodeFromString(strings.NewReader(html), strings.NewReader(""))
		box := nodesToBoxes(sn)
		if !(len(box.children) == 1 && box.children[0].tag() == "p") {
			t.Errorf("should have one <p>")
		}
	})

	t.Run("", func(t *testing.T) {
		html := "<p>1</p><p>2</p>"
		sn := makeStyledNodeFromString(strings.NewReader(html), strings.NewReader(""))
		box := nodesToBoxes(sn)
		if !(len(box.children) == 2 && box.children[0].tag() == "p" && box.children[1].tag() == "p") {
			t.Errorf("should have two <p>")
		}
	})

	t.Run("", func(t *testing.T) {
		var html = `<P>Several <EM>emphasized words</EM> appear
<STRONG>in this</STRONG> sentence, dear.</P>`
		sn := makeStyledNodeFromString(strings.NewReader(html), strings.NewReader(""))
		box := nodesToBoxes(sn)
		if !(len(box.children) == 1 && box.children[0].tag() == "p") {
			t.Errorf("should have one <p>")
		}
	})
}
