package main

type propertyMap map[string]Value

type styledNode struct {
	node            *Node
	specifiedValues []*Declarator
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
		if class != *node.Class() {
			return false
		}
	}

	return false
}

func matchRule(node *Node, rule *Rule) []*Declarator {
	declarators := []*Declarator{}
	for _, sel := range rule.selectors {
		if matches(node, sel) {
			declarators = append(declarators, rule.declarators...)
		}
	}
	return declarators
}

func matchRules(node *Node, style *Stylesheet) []*Declarator {
	declarators := []*Declarator{}
	for _, rule := range style.rules {
		declarators = append(declarators, matchRule(node, &rule)...)
	}
	return declarators
}

func styleTree(node *Node, style *Stylesheet) *styledNode {
	children := []*styledNode{}
	for _, child := range node.Children {
		children = append(children, styleTree(child, style))
	}
	return &styledNode{
		node:            node,
		specifiedValues: matchRules(node, style),
		children:        children,
	}
}
