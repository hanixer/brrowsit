package main

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
