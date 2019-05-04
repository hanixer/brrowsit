package main

import (
	"fmt"
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
.a { background-color: #ff0000; }
.b { background-color: #ffa500; }
.c { background-color: #ffff00; }
.d { background-color: #008000; }
.e { background-color: #0000ff; }
.f { background-color: #4b0082; }
.g { background-color: #800080; }`

var html3 = `
<div class="c"><div class="b"></div></div>`

var css3 = `div { display: block; padding: 12px; }
.b { background-color: #ffa500; }
.c { background-color: #ffff00; }`

var html4 = `
<div class="a">
    <div class="b">
        <div class="c">
        </div>
    </div>
    <div class="d">
        <div class="e">
        </div>
    </div>
</div>`

var css4 = `div { display: block; padding: 12px; }
.a { background-color: #ff0000; }
.b {
	background-color: #ffa500;
	margin-left: 30px; }
.c { background-color: #ffff00; }
.d { background-color: #008000; }
.e { background-color: #0000ff; }
.f { background-color: #4b0082; }
.g { background-color: #800080; }`

var html5 = `<div>simple text...</div>`
var css5 = `div { display: block; padding: 12px; }`

var layout = newColoredBox(rect{20, 20, 300, 200}, red, []*layoutBox{
	newColoredBox(rect{100, 100, 50, 40}, green, nil),
	newColoredBox(rect{100, 200, 10, 10}, green, nil),
	newColoredBox(rect{100, 300, 10, 10}, green, nil),
})

func drawHTMLAndCSS(htmlReader io.Reader, cssReader io.Reader, width int, height int) *image.RGBA {
	n, err := parseHTML(htmlReader)
	if err != nil {
		log.Fatalln("HTML ERROR.", err)
	}
	s, err := ParseStylesheet(cssReader)
	if err != nil {
		log.Fatalln("CSS ERROR.", err)
	}
	st := styleTree(n, s)
	r := nodesToBoxes(st)
	return layoutAndDraw(r, width, height)
}

func main() {
	img := drawHTMLAndCSS(strings.NewReader(html2), strings.NewReader(css2), 600, 400)
	file, err := os.Create("trash.png")
	fmt.Println("file error", err)
	png.Encode(file, img)
}
