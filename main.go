package main

import (
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"strings"
)

var exampleHanded = NewElementNode("html", nil, []*Node{
	NewElementNode("body", nil, []*Node{
		NewElementNode("h1", nil, []*Node{NewTextNode("Title")}),
		NewElementNode("div", nil, []*Node{
			NewElementNode("p", nil, []*Node{
				NewTextNode("Hello"),
				NewElementNode("em", nil, []*Node{NewTextNode("world")}),
				NewTextNode("!one"),
			}),
		}),
		NewElementNode("h1", nil, []*Node{NewTextNode("Title")}),
	}),
})

var example = `<html>
<body>
    <h1>Title</h1>
    <div id="main" class="test">
        <p>Hello <em>world</em>!one</p>
    </div>
</body>
</html>`

var html = `<div>
	<p></p>
	<div id="inner"></div>
</div>`

var css = `#inner {
	background-color: red
}`

var html2 = `<div class="a">
<div class="b">
  <div class="c">
	<div class="d">
	  <div class="e">
		<div class="f">
		  <div class="g">
		  </div>
		</div>
	  </div>
	</div>
  </div>
</div>
</div>`

var css2 = `div { display: block; padding: 12px; }
.a { background: #ff0000; }
.b { background: #ffa500; }
.c { background: #ffff00; }
.d { background: #008000; }
.e { background: #0000ff; }
.f { background: #4b0082; }
.g { background: #800080; }`

var html3 = `
<div class="b"></div>
<div class="c"></div>`

var css3 = `div { display: block; padding: 12px; }
.b { background: #ffa500; }
.c { background: #ffff00; }`

var layout = newColoredBox(rect{20, 20, 300, 200}, red, []*layoutBox{
	newColoredBox(rect{100, 100, 50, 40}, green, nil),
	newColoredBox(rect{100, 200, 10, 10}, green, nil),
	newColoredBox(rect{100, 300, 10, 10}, green, nil),
})

func drawHTMLAndCSS(htmlReader io.Reader, cssReader io.Reader, width int, height int) *image.RGBA {
	n := ParseHtml(htmlReader)
	s, err := ParseStylesheet(cssReader)
	if err != nil {
		log.Fatalln("CSS ERROR.", err)
	}
	st := styleTree(n, s)
	r := nodesToBoxes(st)
	return layoutAndDraw(r, width, height)
}

func main() {
	img := drawHTMLAndCSS(strings.NewReader(html3), strings.NewReader(css3), 600, 400)
	file, _ := os.Create("trash.png")
	png.Encode(file, img)
}
