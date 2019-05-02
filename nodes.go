package main

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

// GetID returns ID if it's present
func (n *Node) GetID() *string {
	id, ok := n.Attributes["id"]
	if ok {
		return &id
	}

	return nil
}

// Class return class
func (n *Node) Class() *string {
	id, ok := n.Attributes["class"]
	if ok {
		return &id
	}

	return nil
}

func (n *Node) TagName() string {
	return n.Data
}
