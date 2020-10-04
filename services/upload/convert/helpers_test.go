package main

import (
	"reflect"
	"testing"
)

func TestSplitAtCommas(t *testing.T) {

	t.Run("split by comma without qoutes", func(t *testing.T) {
		row := "hello,world"
		got := splitAtCommas(&row)
		want := []string{"hello", "world"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("split by comma with qoutes", func(t *testing.T) {
		row := "\"hello, world\",yo"
		got := splitAtCommas(&row)
		want := []string{"\"hello, world\"", "yo"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestTrimQuotes(t *testing.T) {
	t.Run("trim existing quotes", func(t *testing.T) {
		row := "\"hello, world\""
		got := trimQuotes(row)
		want := "hello, world"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("don't trim nonexisting quotes", func(t *testing.T) {
		row := "hello, world"
		got := trimQuotes(row)
		want := "hello, world"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
