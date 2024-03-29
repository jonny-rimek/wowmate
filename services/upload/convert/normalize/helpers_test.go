package normalize

import (
	"reflect"
	"testing"
)

func TestSplitAtCommas(t *testing.T) {

	t.Run("split by comma without quotes", func(t *testing.T) {
		row := "hello,world"
		got := splitAtCommas(&row)
		want := []string{"hello", "world"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("split by comma with quotes", func(t *testing.T) {
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

	t.Run("don't trim nonexistent quotes", func(t *testing.T) {
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

func TestAtoi64(t *testing.T) {
	t.Run("convert number as string to int64", func(t *testing.T) {
		got, err := Atoi64("32")
		want := int64(32)

		if err != nil {
			t.Fatal("shouldn't fail")
		}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("fail if text is passed", func(t *testing.T) {
		got, err := Atoi64("hello, world")
		want := int64(0)

		if err == nil {
			t.Fatal("must fail")
		}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("fail if the number is too big", func(t *testing.T) {
		got, err := Atoi64("9223372036854775808")
		want := int64(0)

		if err == nil {
			t.Fatal("must fail")
		}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("convert max int64 number", func(t *testing.T) {
		got, err := Atoi64("9223372036854775807")
		want := int64(9223372036854775807)

		if err != nil {
			t.Fatal("mustn't fail")
		}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("convert max negative int64 number", func(t *testing.T) {
		got, err := Atoi64("-9223372036854775807")
		want := int64(-9223372036854775807)

		if err != nil {
			t.Fatal("mustn't fail")
		}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func Test_splitString(t *testing.T) {
	input := "4/24 10:42:30.561  COMBAT_LOG_VERSION"
	want := []string{"4/24 10:42:30.561", "COMBAT_LOG_VERSION"}

	got := splitString(input, "  ")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q want %q", got, want)
	}
}
