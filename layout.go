package main

// LayoutObject
// - LayoutText
// - LayoutBox
// -- LayoutBlock
// -- LayoutInline
// -- LayoutLineBox

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"
)

type dimensions struct {
	content rect

	padding edgeSizes
	margin  edgeSizes
	border  edgeSizes
}

type rect struct {
	x, y, width, height float32
}

type edgeSizes struct {
	left, right, top, bottom float32
}

type layoutBox struct {
	dimensions dimensions
	boxType    boxType
	styledNode *styledNode
	children   []*layoutBox
}

type boxType int

const (
	blockBox boxType = iota
	inlineBox
	textBox
	anonymousBox
	rootBox
)

type bounds struct {
	width, height int
}

var blockTags = []string{
	"div", "p",
}

func (r rect) min() image.Point {
	return image.Point{int(r.x), int(r.y)}
}

func (r rect) expandedBy(s edgeSizes) rect {
	return rect{
		x:      r.x - s.left,
		y:      r.y - s.top,
		width:  r.width + s.left + s.right,
		height: r.height + s.top + s.bottom,
	}
}

func (d dimensions) paddingBox() rect {
	return d.content.expandedBy(d.padding)
}

func (d dimensions) borderBox() rect {
	return d.paddingBox().expandedBy(d.border)
}

func (d dimensions) marginBox() rect {
	return d.borderBox().expandedBy(d.margin)
}

func isJustSpaces(n *Node) bool {
	return len(strings.TrimSpace(n.Data)) == 0
}

func makeAnonymousBlockBox() *layoutBox {
	box := new(layoutBox)
	box.boxType = blockBox
	box.children = []*layoutBox{}
	box.styledNode = nil
	return box
}

func nodesToBoxes(node *styledNode) *layoutBox {
	childBoxes := []*layoutBox{}
	isOnlyInline := true

	for _, child := range node.children {
		childBox := nodesToBoxes(child)
		if childBox == nil {
			continue
		}

		if isOnlyInline && childBox.isBlock() {
			isOnlyInline = false
			if len(childBoxes) > 0 {
				anonymous := makeAnonymousBlockBox()
				anonymous.children = childBoxes
				childBoxes = []*layoutBox{anonymous}
			}
			childBoxes = append(childBoxes, childBox)
		} else if !isOnlyInline && !childBox.isBlock() {
			// the child should be added to anonymous box
			last := childBoxes[len(childBoxes)-1]
			if last.isAnonymous() {
				last.children = append(last.children, childBox)
			} else {
				anonymous := makeAnonymousBlockBox()
				anonymous.children = []*layoutBox{childBox}
				childBoxes = append(childBoxes, anonymous)
			}
		} else {
			childBoxes = append(childBoxes, childBox)
		}
	}

	box := new(layoutBox)
	box.children = childBoxes
	box.styledNode = node

	if box.styledNode.node.NodeType == RootNode {
		box.boxType = rootBox
	} else if box.styledNode.node.NodeType == TextNode {
		box.boxType = textBox
		if isJustSpaces(node.node) {
			box = nil
		}
	} else if isBlockElement(node) {
		box.boxType = blockBox
	} else {
		box.boxType = inlineBox
	}

	return box
}

func (box *layoutBox) layoutRoot(width, height int) {
	box.dimensions = dimensions{}
	box.dimensions.content.width = float32(width)

	box.layoutChildren()
}

func newLineBox(x, y, width float32) *layoutBox {
	b := new(layoutBox)
	b.boxType = anonymousBox
	b.children = []*layoutBox{}
	b.dimensions.content.x = x
	b.dimensions.content.y = y
	b.dimensions.content.width = width
	return b
}

func (box *layoutBox) nextXPos() float32 {
	if len(box.children) == 0 {
		return box.dimensions.content.x
	}
	last := box.children[len(box.children)-1]
	mb := last.dimensions.marginBox()
	return mb.x + mb.width
}

func (box *layoutBox) canAppendToLine(inlineBox *layoutBox) bool {
	x := box.nextXPos()
	mb := inlineBox.dimensions.marginBox()
	return x+mb.width < box.dimensions.content.width
}

func (box *layoutBox) appendToLine(inlineBox *layoutBox) {
	x := box.nextXPos()
	inlineBox.dimensions.content.x = x
	inlineBox.dimensions.content.y = box.dimensions.content.y
	box.children = append(box.children, inlineBox)
}

func (box *layoutBox) calculateLineHeight() {
	var h float32
	for _, child := range box.children {
		if child.dimensions.content.height > h {
			h = child.dimensions.content.height
		}
	}
}

func (box *layoutBox) layoutChildren() {
	newChildren := []*layoutBox{}
	var lineBox *layoutBox
	x := box.dimensions.content.x
	y := box.dimensions.content.y
	width := box.dimensions.content.width
	for _, child := range box.children {
		if child.boxType != blockBox {
			if lineBox == nil {
				lineBox = newLineBox(x, y, width)
			}
			if !lineBox.canAppendToLine(child) {
				newChildren = box.appendLine(newChildren, lineBox)
				lineBox = newLineBox(x, y, width)
			}
			child.layout(box.dimensions)
			lineBox.appendToLine(child)
		} else {
			child.layout(box.dimensions)
			newChildren = append(newChildren, child)
			box.dimensions.content.height += child.dimensions.marginBox().height
		}
	}
	if lineBox != nil {
		newChildren = box.appendLine(newChildren, lineBox)
	}
	box.children = newChildren
}

func (box *layoutBox) isBlock() bool {
	return box.boxType == blockBox
}

func (box *layoutBox) isAnonymous() bool {
	return box.boxType != textBox && box.styledNode == nil
}

func (box *layoutBox) tag() string {
	if (box.boxType == blockBox || box.boxType == inlineBox) && box.styledNode != nil && box.styledNode.node != nil {
		return box.styledNode.node.TagName()
	}
	return ""
}

func (box *layoutBox) appendLine(newChildren []*layoutBox, accumulator *layoutBox) []*layoutBox {
	newChildren = append(newChildren, accumulator)
	accumulator.calculateLineHeight()
	box.dimensions.content.height += accumulator.dimensions.marginBox().height
	return newChildren
}

func (box *layoutBox) layout(containingBlock dimensions) {
	if box.boxType == blockBox {
		box.calculateWidth(containingBlock)
		box.calculatePosition(containingBlock)
	}

	box.layoutChildren()

	if box.boxType != blockBox {
		for _, child := range box.children {
			box.dimensions.content.width += child.dimensions.content.width
		}
		if box.boxType == textBox {
			box.dimensions.content.width = float32(getStringWidth(box.styledNode.node.Data))
			box.dimensions.content.height = float32(getFontHeight())
		}
	}
}

// a lot more simple than specified https://www.w3.org/TR/CSS2/visudet.html#Computing_widths_and_margins
func (box *layoutBox) calculateWidth(containingBlock dimensions) {
	node := box.styledNode
	if node == nil {
		return
	}

	auto := Value{keyword: "auto"}
	width := node.lookupOr("width", auto)
	zero := Value{length: 0, valueType: Length}
	d := &box.dimensions

	marginLeft := node.lookupDouble("margin-left", "margin", zero)
	marginRight := node.lookupDouble("margin-right", "margin", zero)
	borderLeft := node.lookupDouble("border-left-width", "border-width", zero)
	borderRight := node.lookupDouble("border-right-width", "border-width", zero)
	paddingLeft := node.lookupDouble("padding-left", "padding", zero)
	paddingRight := node.lookupDouble("padding-right", "padding", zero)

	d.margin.left = marginLeft.toPx()
	d.margin.right = marginRight.toPx()
	d.border.left = borderLeft.toPx()
	d.border.right = borderRight.toPx()
	d.padding.left = paddingLeft.toPx()
	d.padding.right = paddingRight.toPx()

	if width == auto {
		sum := d.margin.left + d.margin.right
		sum += d.border.left + d.border.right
		sum += d.padding.left + d.padding.right

		d.content.width = containingBlock.content.width - sum
	} else {
		d.content.width = width.toPx()
	}
}

func (box *layoutBox) calculatePosition(containingBlock dimensions) {
	node := box.styledNode
	if node == nil {
		return
	}

	zero := Value{length: 0, valueType: Length}
	d := &box.dimensions

	marginTop := node.lookupDouble("margin-top", "margin", zero)
	marginBottom := node.lookupDouble("margin-bottom", "margin", zero)
	borderTop := node.lookupDouble("border-top-width", "border-width", zero)
	borderBottom := node.lookupDouble("border-bottom-width", "border-width", zero)
	paddingTop := node.lookupDouble("padding-top", "padding", zero)
	paddingBottom := node.lookupDouble("padding-bottom", "padding", zero)

	d.margin.top = marginTop.toPx()
	d.margin.bottom = marginBottom.toPx()
	d.border.top = borderTop.toPx()
	d.border.bottom = borderBottom.toPx()
	d.padding.top = paddingTop.toPx()
	d.padding.bottom = paddingBottom.toPx()

	d.content.x = containingBlock.content.x + d.margin.left + d.border.left + d.padding.left

	d.content.y = containingBlock.content.y + containingBlock.content.height
	d.content.y += d.margin.top + d.border.top + d.padding.top
}

var red = color.RGBA{255, 0, 0, 255}
var green = color.RGBA{0, 255, 0, 255}

func newColoredBox(r rect, c color.RGBA, children []*layoutBox) *layoutBox {
	node := &styledNode{specifiedValues: propertyMap{"background-color": Value{color: c}}}
	return &layoutBox{
		dimensions: dimensions{content: r},
		styledNode: node,
		children:   children,
	}
}

// To layout a node:
// If it has block display:
//  Start alignment from left border
//	Layout recursively
//	We receive a box
//	We need to know the height of this box. This available in the box
func isBlockElement(node *styledNode) bool {
	v, ok := node.specifiedValues["display"]
	if ok {
		return v.keyword == "block"
	}

	for _, v := range blockTags {
		if v == node.node.TagName() {
			return true
		}
	}
	return false
}

// PrintLayoutTree prints the tree
func PrintLayoutTree(box *layoutBox, w io.Writer) {
	printLayoutTree(box, w, 0)
}

func printLayoutTree(box *layoutBox, w io.Writer, nesting int) {
	fmt.Fprint(w, box.boxType)
	if box.tag() != "" {
		fmt.Fprintf(w, " <%s> ", box.tag())
	}
	if box.boxType == textBox {
		fmt.Fprintf(w, " %q", box.styledNode.node.Data)
	}
	for _, c := range box.children {
		printNesting(w, nesting+1)
		printLayoutTree(c, w, nesting+1)
	}
}

func (box *layoutBox) String() string {
	build := new(strings.Builder)
	PrintLayoutTree(box, build)
	return build.String()
}
