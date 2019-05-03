package main

import (
	"sort"
)

type propertyMap map[string]Value

type styledNode struct {
	node            *Node
	specifiedValues propertyMap
	children        []*styledNode
}

type displayType int

const (
	inline displayType = iota
	block
	none
)

func (n *styledNode) displayType() displayType {
	if n.node.NodeType == ElementNode {
		v, ok := n.lookup("display")
		if ok {
			switch v.keyword {
			case "block":
				return block
			case "inline":
				return inline
			case "none":
				return none
			}
		}
		return block
	}
	return inline
}

func matches(node *Node, selector *Selector) bool {
	if node.NodeType == TextNode {
		return false
	}

	if selector.tagName != nil && *selector.tagName != node.TagName() {
		return false
	}

	nodeID := node.GetID()
	if selector.id != nil && nodeID != nil && *selector.id != *nodeID {
		return false
	}

	if node.Class() == nil {
		return true
	}

	for _, class := range selector.class {
		if class != *node.Class() {
			return false
		}
	}

	return true
}

func matchRule(node *Node, rule *Rule) bool {
	for _, sel := range rule.selectors {
		if matches(node, sel) {
			return true
		}
	}
	return false
}

func matchRules(node *Node, style *Stylesheet) propertyMap {
	pmap := make(propertyMap)
	if node.NodeType != ElementNode {
		return pmap
	}
	for _, rule := range style.rules {
		if matchRule(node, rule) {
			for _, decl := range rule.declarators {
				pmap[decl.name] = decl.value
			}
		}
	}
	return pmap
}

func styleTree(node *Node, style *Stylesheet) *styledNode {
	children := []*styledNode{}
	sort.Slice(style.rules, func(i, j int) bool {
		return compareSpecificity(style.rules[i].selectors, style.rules[j].selectors) < 0
	})
	for _, child := range node.Children {
		children = append(children, styleTree(child, style))
	}
	return &styledNode{
		node:            node,
		specifiedValues: matchRules(node, style),
		children:        children,
	}
}

func (n *styledNode) lookup(k string) (Value, bool) {
	v, ok := n.specifiedValues[k]
	return v, ok
}

func (n *styledNode) lookupOr(k string, elseVal Value) Value {
	v, ok := n.specifiedValues[k]
	if ok {
		return v
	}
	return elseVal
}

func (n *styledNode) lookupDouble(k1, k2 string, elseVal Value) Value {
	v, ok := n.lookup(k1)
	if ok {
		return v
	}
	return n.lookupOr(k2, elseVal)
}
