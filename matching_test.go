package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_styleTree(t *testing.T) {
	type args struct {
		node  *Node
		style *Stylesheet
	}
	tests := []struct {
		name string
		args args
		want *styledNode
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := styleTree(tt.args.node, tt.args.style); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("styleTree() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("", func(t *testing.T) {
		html := `<p>things</p>`
		css := `p {font-weight: bold}`
		node := ParseHtml(strings.NewReader(html))
		style, err := ParseStylesheet(strings.NewReader(css))
		if err != nil {
			t.Errorf("css parse error")
		}
		styled := styleTree(node, style)
		fmt.Println(styled)
	})
}
