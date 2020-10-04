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

func TestAtoi32(t *testing.T) {
	t.Run("convert number as string to int32", func(t *testing.T) {
		got, err := Atoi32("32")
		want := int32(32)

		if err != nil {
			t.Fatal("shouldn't fail")
		}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("fail if text is passed", func(t *testing.T) {
		got, err := Atoi32("hello, world")
		want := int32(0)

		if err == nil {
			t.Fatal("must fail")
		}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("fail if the number is too big", func(t *testing.T) {
		got, err := Atoi32("2147483648")
		want := int32(0)

		if err == nil {
			t.Fatal("must fail")
		}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("convert max int32 number", func(t *testing.T) {
		got, err := Atoi32("2147483647")
		want := int32(2147483647)

		if err != nil {
			t.Fatal("mustn't fail")
		}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("convert max negative int32 number", func(t *testing.T) {
		got, err := Atoi32("-2147483647")
		want := int32(-2147483647)

		if err != nil {
			t.Fatal("mustn't fail")
		}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
