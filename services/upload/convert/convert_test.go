package main

import (
	"testing"
)

func Test_uploadUUID(t *testing.T) {
	t.Run("get upload uuid from s3 key with .txt", func(t *testing.T) {
		key := "upload/2020/10/7/23/e894fa50-c4ba-42e3-ae87-b9076630f2b6.txt"
		got := uploadUUID(key)
		want := "e894fa50-c4ba-42e3-ae87-b9076630f2b6"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("get upload uuid from s3 key with .txt.gz", func(t *testing.T) {
		key := "upload/2020/10/7/23/e894fa50-c4ba-42e3-ae87-b9076630f2b6.txt.gz"
		got := uploadUUID(key)
		want := "e894fa50-c4ba-42e3-ae87-b9076630f2b6"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("get upload uuid from s3 key with .zip", func(t *testing.T) {
		key := "upload/2020/10/7/23/e894fa50-c4ba-42e3-ae87-b9076630f2b6.zip"
		got := uploadUUID(key)
		want := "e894fa50-c4ba-42e3-ae87-b9076630f2b6"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
