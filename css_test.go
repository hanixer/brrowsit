package main

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func mr(s string) *bufio.Reader { return bufio.NewReader(strings.NewReader(s)) }

func Test_parseSelector(t *testing.T) {
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Selector
		wantErr bool
	}{
		// TODO: Add test cases.
		// {"", args{mr("")}, Selector{tagName: ""}, false},
	}
	t.Run("", func(t *testing.T) {
		r := mr("")
		got, err := parseSelector(r)

		if got != nil || err != nil {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("a")
		got, err := parseSelector(r)

		if err != nil || got == nil || got.tagName == nil || *got.tagName != "a" {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("a#id.class1.class2")
		got, err := parseSelector(r)

		if err != nil || *got.tagName != "a" || *got.id != "id" || len(got.class) < 2 ||
			got.class[0] != "class1" || got.class[1] != "class2" {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("#id1")
		got, err := parseSelector(r)

		if err != nil || *got.id != "id1" {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("a  ,")
		got, err := parseSelector(r)

		if err != nil || got == nil || got.tagName == nil || *got.tagName != "a" {
			t.Errorf("%v ::: %v", got, err)
		}
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSelector(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSelector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSelector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSelectors(t *testing.T) {
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []Selector
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSelectors(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSelectors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSelectors() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("", func(t *testing.T) {
		r := mr("")
		got, err := parseSelectors(r)

		if err != nil || len(got) != 0 {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("a,b,                               c         ")
		got, err := parseSelectors(r)

		if err != nil || len(got) != 3 {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("a,b,                               c    {     ")
		got, err := parseSelectors(r)

		if err != nil || len(got) != 3 {
			t.Error(got, err)
		}
	})
}

func Test_parseDeclarator(t *testing.T) {
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Declarator
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDeclarator(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDeclarator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDeclarator() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("", func(t *testing.T) {
		r := mr("")
		got, err := parseDeclarator(r)

		if !(err == nil && got == nil) {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("#")
		got, err := parseDeclarator(r)

		if !(err == nil && got == nil) {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("a: b")
		got, err := parseDeclarator(r)

		if !(err == nil && got != nil && got.name == "a" && got.value.keyword == "b" && got.valueType == Keyword) {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("a: 123px")
		got, err := parseDeclarator(r)

		if !(err == nil && got != nil && got.name == "a" && got.value.length == 123 && got.valueType == Length) {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("a: #010203")
		got, err := parseDeclarator(r)
		if err == nil && got != nil && got.name == "a" {
			c := got.value.color
			if c.r == 1 && c.g == 2 && c.b == 3 {
			} else {
				t.Error(got, err)
			}
		} else {
			t.Error(got, err)
		}
	})
}

func Test_parseDeclarators(t *testing.T) {
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []Declarator
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDeclarators(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDeclarators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDeclarators() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("", func(t *testing.T) {
		r := mr("{}")
		got, err := parseDeclarators(r)

		if !(err == nil && got != nil && len(got) == 0) {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("{a:b;c:d}")
		got, err := parseDeclarators(r)

		if !(err == nil && got != nil && len(got) == 2) {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("{a:b;c:d;}")
		got, err := parseDeclarators(r)

		if !(err == nil && got != nil && len(got) == 2) {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr("{;a:b;c:d;}")
		got, err := parseDeclarators(r)

		if !(err == nil && got != nil && len(got) == 2) {
			t.Error(got, err)
		}
	})
}

func TestParseStylesheet(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Stylesheet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStylesheet(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStylesheet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseStylesheet() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("", func(t *testing.T) {
		r := mr(``)
		got, err := ParseStylesheet(r)

		if !(err == nil && got != nil && len(got.rules) == 0) {
			t.Error(got, err)
		}
	})
	t.Run("", func(t *testing.T) {
		r := mr(`h1, h2, h3 { margin: auto; color: #cc0000; }
		div.note { margin-bottom: 20px; padding: 10px; }
		#answer { display: none; }`)
		got, err := ParseStylesheet(r)

		if !(err == nil && got != nil && len(got.rules) == 3) {
			t.Error(got, err)
		}
	})

}

func Test_compareSpecificity(t *testing.T) {
	type args struct {
		sels1 []*Selector
		sels2 []*Selector
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareSpecificity(tt.args.sels1, tt.args.sels2); got != tt.want {
				t.Errorf("compareSpecificity() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("", func(t *testing.T) {
		ss1, _ := ParseStylesheet(mr(`.one, .two, h1 {}`))
		ss2, _ := ParseStylesheet(mr(`p, .two, .three {}`))
		fmt.Println(ss1)

		if !(compareSpecificity(ss1.rules[0].selectors, ss2.rules[0].selectors) == 0) {
			t.Error(ss1, ss2)
		}
	})
}
