package main

import "testing"

func Test_getStringWidth(t *testing.T) {
	t.Run("", func(t *testing.T) {
		s := "Some test string."
		w1 := getStringWidth(s)
		w2 := getStringWidth(s + s)
		if w1*2 != w2 {
			t.Errorf("%d * 2 != %d", w1, w2)
		}
	})
}
