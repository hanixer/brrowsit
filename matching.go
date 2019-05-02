package main

import (
	"fmt"
	"sort"
)

type propertyMap map[string]Value

type styledNode struct {
	node            *Node
	specifiedValues propertyMap
	children        []*styledNode
}

func matches(node *Node, selector *Selector) bool {
	if node.NodeType == TextNode {
		return false
	}

	if selector.tagName != nil && *selector.tagName != node.TagName() {
		return false
	}

	if selector.id != nil && *selector.id != *node.GetID() {
		return false
	}

	for _, class := range selector.class {
		fmt.Println("class is", *node.Class())
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
