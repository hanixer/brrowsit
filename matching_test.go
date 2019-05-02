package main

import (
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
		html := `<p>things</p><div id="uniq">nothings</div>`
		css := `p {font-weigth: bold} div {color: #001122}`
		node := ParseHtml(strings.NewReader(html))
		style, err := ParseStylesheet(strings.NewReader(css))
		if err != nil {
			t.Errorf("css parse error")
		}
		styled := styleTree(node, style)
		if !(styled.node.NodeType == RootNode) {
			t.Errorf("show be root")
		}
		if !(len(styled.children) == 2) {
			t.Error("should be 2 children")
		}
		ch1 := styled.children[0]
		if !(ch1.node.TagName() == "p" && pmapContainsKey(ch1.specifiedValues, "font-weigth")) {
			t.Error("wrong 1st node")
		}
		ch1 = styled.children[1]
		if !(ch1.node.TagName() == "div" && pmapContainsKey(ch1.specifiedValues, "color") &&
			pmapHasType(ch1.specifiedValues, "color", ColorValue)) {
			t.Error("wrong 2nd node")
		}
	})
	t.Run("test specifity", func(t *testing.T) {
		html := `<p id="koin">things</p>`
		css := `p {font-weight: bold} #koin {color: #001122; font-weight: nono}`
		node := ParseHtml(strings.NewReader(html))
		style, err := ParseStylesheet(strings.NewReader(css))
		if err != nil {
			t.Errorf("css parse error")
		}
		styled := styleTree(node, style)
		if !(styled.node.NodeType == RootNode) {
			t.Errorf("show be root")
		}
		if !(len(styled.children) == 1) {
			t.Error("should be 1 children")
		}
		ch1 := styled.children[0]
		if !(ch1.node.TagName() == "p" && pmapContainsKey(ch1.specifiedValues, "font-weight") &&
			ch1.specifiedValues["font-weight"].keyword == "nono") {
			t.Error("Wrong properies, got", ch1.specifiedValues["font-weight"].keyword)
		}
	})
}

func pmapContainsKey(decls propertyMap, k string) bool {
	_, ok := decls[k]
	return ok
}

func pmapHasType(decls propertyMap, k string, t valueType) bool {
	v, ok := decls[k]
	if ok {
		return v.valueType == t
	}
	return false
}
