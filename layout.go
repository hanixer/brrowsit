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

type displayType int

const (
	inline displayType = iota
	block
	none
)

type bounds struct {
	width, height int
}

func layoutBoxes(bounds bounds, rootBox *layoutBox) {

}

var red = color.RGBA{255, 0, 0, 255}
var green = color.RGBA{0, 255, 0, 255}

func newColoredBox(r rect, c color.RGBA, children []*layoutBox) *layoutBox {
	node := &styledNode{specifiedValues: propertyMap{"color": Value{color: c}}}
	return &layoutBox{
		dimensions: dimensions{content: r},
		styledNode: node,
		children:   children,
	}
}
