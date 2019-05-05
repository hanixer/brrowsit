package main

import (
	"fmt"
	"image"
	"image/color"
)

type drawCommand interface {
	draw(img *image.RGBA)
}

type drawRect struct {
	color color.RGBA
	rect  rect
}

type drawText struct {
	s  string
	pt image.Point
}

func (d *drawRect) draw(img *image.RGBA) {
	fmt.Printf("x=%v, y=%v, w=%v, h=%v\n", d.rect.x, d.rect.y, d.rect.width, d.rect.height)
	for x := d.rect.x; x < d.rect.x+d.rect.width; x++ {
		for y := d.rect.y; y < d.rect.y+d.rect.height; y++ {
			img.Set(int(x), int(y), d.color)
		}
	}
}

func (d *drawText) draw(img *image.RGBA) {
	drawString(d.s, img, d.pt)
}

func mergeLists(l1 []drawCommand, l2 []drawCommand) []drawCommand {
	for _, elem := range l2 {
		l1 = append(l1, elem)
	}
	return l1
}

func makeDisplayList(layout *layoutBox) []drawCommand {
	commands := []drawCommand{}
	if layout.boxType == blockBox {
		v, ok := layout.styledNode.specifiedValues["background-color"]
		if ok {

		}
		d := &drawRect{v.color, layout.dimensions.paddingBox()}
		commands = append(commands, d)
	} else if layout.boxType == textBox {
		d := &drawText{layout.styledNode.node.Data, layout.dimensions.content.min()}
		commands = append(commands, d)
	}

	for _, child := range layout.children {
		commands = mergeLists(commands, makeDisplayList(child))
	}

	return commands
}

func drawDisplayList(img *image.RGBA, commands []drawCommand) {
	for _, comm := range commands {
		comm.draw(img)
	}
}

func layoutAndDraw(rootBox *layoutBox, width int, height int) *image.RGBA {
	rootBox.layoutRoot(width, height)
	commands := makeDisplayList(rootBox)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	drawDisplayList(img, commands)
	return img
}
