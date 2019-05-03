package main

import "image/color"

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
	blockNode boxType = iota
	inlineNode
	anonymous
)

type bounds struct {
	width, height int
}

var blockTags = []string{
	"div", "p",
}

func nodesToBoxes(node *styledNode) *layoutBox {
	childBoxes := []*layoutBox{}

	for _, child := range node.children {
		childBoxes = append(childBoxes, nodesToBoxes(child))
	}

	box := new(layoutBox)
	box.children = childBoxes
	box.boxType = blockNode
	box.styledNode = node

	// switch node.displayType() {
	// case block:
	// 	// TODO
	// case inline:
	// 	// TODO:
	// case none:
	// 	// nothing TODO
	// }
	return box
}

func (box *layoutBox) layoutRoot(width, height int) {
	d := dimensions{}
	d.content.width = float32(width)
	d.content.height = float32(height)

	for _, child := range box.children {
		child.layout(box.dimensions)

		box.calculateHeight(child)
	}
}

func (box *layoutBox) layout(containingBlock dimensions) {
	box.calculateWidth(containingBlock)
	box.calculatePosition(containingBlock)

	for _, child := range box.children {
		child.layout(box.dimensions)

		box.calculateHeight(child)
	}
}

func (box *layoutBox) calculateHeight(lastChild *layoutBox) {
	d := &box.dimensions
	cd := &lastChild.dimensions
	d.content.height += cd.border.bottom + cd.padding.bottom + cd.margin.bottom + cd.content.height
}

// a lot more simple than specified https://www.w3.org/TR/CSS2/visudet.html#Computing_widths_and_margins
func (box *layoutBox) calculateWidth(containingBlock dimensions) {
	node := box.styledNode
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
