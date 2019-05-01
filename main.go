package main

import (
	"fmt"
	"strings"
)

type NodeType int

const (
	TextNode NodeType = iota
	ElementNode
)

type Node struct {
	NodeType   NodeType
	Data       string
	Children   []*Node
	Attributes map[string]string
}

func NewTextNode(s string) *Node {
	return &Node{TextNode, s, []*Node{}, make(map[string]string)}
}

func NewElementNode(name string, attrs map[string]string, ch []*Node) *Node {
	return &Node{ElementNode, name, ch, attrs}
}

var exampleHanded = NewElementNode("html", nil, []*Node{
	NewElementNode("body", nil, []*Node{
		NewElementNode("h1", nil, []*Node{NewTextNode("Title")}),
		NewElementNode("div", nil, []*Node{
			NewElementNode("p", nil, []*Node{
				NewTextNode("Hello"),
				NewElementNode("em", nil, []*Node{NewTextNode("world")}),
				NewTextNode("!one"),
			}),
		}),
		NewElementNode("h1", nil, []*Node{NewTextNode("Title")}),
	}),
})

var example = `<html>
<body>
    <h1>Title</h1>
    <div id="main" class="test">thing
        <p>Hello <em>world</em>!one</p>
    </div>
</body>
</html>`

func main() {
	p := NewParser(strings.NewReader(example))
	fmt.Println(p.Parse())
}
