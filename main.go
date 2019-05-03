package main

import (
	"image"
	"image/png"
	"os"
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

var layout = newColoredBox(rect{20, 20, 300, 200}, red, []*layoutBox{
	newColoredBox(rect{100, 100, 50, 40}, green, nil),
	newColoredBox(rect{100, 200, 10, 10}, green, nil),
	newColoredBox(rect{100, 300, 10, 10}, green, nil),
})

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	commands := makeDisplayList(layout)
	drawDisplayList(img, commands)
	file, _ := os.Create("trash.png")
	png.Encode(file, img)
}
