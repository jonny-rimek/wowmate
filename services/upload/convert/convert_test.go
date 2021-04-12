package main

import (
	"testing"
)

func Test_uploadUUID(t *testing.T) {
	assertCorrectUUID := func(t testing.TB, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}

	t.Run("get upload uuid from s3 key with .txt", func(t *testing.T) {
		key := "upload/2020/10/7/23/e894fa50-c4ba-42e3-ae87-b9076630f2b6.txt"
		got, _ := uploadUUID(key)
		want := "e894fa50-c4ba-42e3-ae87-b9076630f2b6"

		assertCorrectUUID(t, got, want)
	})

	t.Run("get upload uuid from s3 key with .txt.gz", func(t *testing.T) {
		key := "upload/2020/10/7/23/e894fa50-c4ba-42e3-ae87-b9076630f2b6.txt.gz"
		got, _ := uploadUUID(key)
		want := "e894fa50-c4ba-42e3-ae87-b9076630f2b6"

		assertCorrectUUID(t, got, want)
	})

	t.Run("get upload uuid from s3 key with .zip", func(t *testing.T) {
		key := "upload/2020/10/7/23/e894fa50-c4ba-42e3-ae87-b9076630f2b6.zip"
		got, _ := uploadUUID(key)
		want := "e894fa50-c4ba-42e3-ae87-b9076630f2b6"

		assertCorrectUUID(t, got, want)
	})

	t.Run("return error if input is empty", func(t *testing.T) {
		// refactor this test, there should be a better way to test correct expected error
		key := ""
		_, got := uploadUUID(key)
		want := "input can't be empty"

		if got.Error() != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("input has the wrong length", func(t *testing.T) {
		// refactor this test, there should be a better way to test correct expected error
		key := "upload/2020/10/7/e894fa50-c4ba-42e3-ae87-b9076630f2b6.zip" // the hour part is missing
		_, got := uploadUUID(key)
		want := "input has the wrong length, got 5 want 6"

		if got.Error() != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
