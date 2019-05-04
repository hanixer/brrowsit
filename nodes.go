package main

import (
	"fmt"
	"io"
	"strings"
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

// TagName returns tag of element node, empty if not element node
func (n *Node) TagName() string {
	return n.Data
}

// PrintNode pretty-prints the node
func PrintNode(node *Node, w io.Writer) {
	printNode(node, w, -1)
}

func printNode(node *Node, w io.Writer, nesting int) {
	switch node.NodeType {
	case TextNode:
		io.WriteString(w, node.Data)
	case ElementNode, RootNode:
		if node.NodeType == ElementNode {
			fmt.Fprintf(w, "<%s", node.Data)
			printAttributes(w, node.Attributes)
			fmt.Fprint(w, ">")
		}
		newLinesNeccessary := len(node.Children) > 0 && node.Children[0].NodeType != TextNode
		for _, child := range node.Children {
			if newLinesNeccessary {
				printNesting(w, nesting+1)
			}
			printNode(child, w, nesting+1)
		}
		if node.NodeType == ElementNode {
			if newLinesNeccessary {
				printNesting(w, nesting)
			}
			fmt.Fprintf(w, "</%s>", node.Data)
		}
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
