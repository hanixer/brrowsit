package main

import (
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

func (d *drawRect) draw(img *image.RGBA) {
	for x := d.rect.x; x < d.rect.x+d.rect.width; x++ {
		for y := d.rect.y; y < d.rect.y+d.rect.height; y++ {
			img.Set(int(x), int(y), d.color)
		}
	}
}

func mergeLists(l1 []drawCommand, l2 []drawCommand) []drawCommand {
	for _, elem := range l2 {
		l1 = append(l1, elem)
	}
	return l1
}

func makeDisplayList(layout *layoutBox) []drawCommand {
	commands := []drawCommand{}
	if layout.boxType != anonymous {
		v, ok := layout.styledNode.specifiedValues["color"]
		if ok {

		}
		d := &drawRect{v.color, layout.dimensions.content}
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
