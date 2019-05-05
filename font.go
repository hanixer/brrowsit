package main

import (
	"image"
	"image/draw"
	"io/ioutil"
	"log"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

var (
	dpi      = 72.0
	fontfile = "fonts/PlayfairDisplay-Black.ttf"
	hinting  = "none"
	size     = 12.0
	spacing  = 1.5
)

var font = initFont(fontfile)

func initFont(fontFile string) *truetype.Font {
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Fatalln("font loading failed.", err)
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatalln("freetype.ParseFont failed.", err)
	}
	return f
	// fg, bg := image.Black, image.White
	// c := freetype.NewContext()
	// c.SetDPI(dpi)
	// c.SetFont(f)
	// c.SetFontSize(size)
	// // c.SetClip(rgba.Bounds())
	// // c.SetDst(rgba)
	// c.SetSrc(fg)
}

func drawString(s string, img draw.Image, pt image.Point) {
	fg, _ := image.Black, image.White
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)
	ptf := freetype.Pt(pt.X, pt.Y)
	_, err := c.DrawString(s, ptf)
	if err != nil {
		log.Println(err)
		return
	}
}

func getStringWidth(s string) int {
	face := truetype.NewFace(font, nil)
	total := 0
	for _, r := range s {
		p, ok := face.GlyphAdvance(r)
		if !ok {
			log.Fatalln("No glyph for", string(r))
		}
		total += p.Round()
	}
	return total
}

func getFontHeight() int {
	face := truetype.NewFace(font, nil)
	return face.Metrics().Height.Round()
}
