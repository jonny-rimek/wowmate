package main

import (
	"testing"
)

func Test_uploadUUID(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "get upload uuid from s3 key with .txt",
			args: args{s: "upload/2020/10/7/23/e894fa50-c4ba-42e3-ae87-b9076630f2b6.txt"},
			want: "e894fa50-c4ba-42e3-ae87-b9076630f2b6",
			wantErr: false,
		},
		{
			name: "get upload uuid from s3 key with .zip",
			args: args{s: "upload/2020/10/7/23/e894fa50-c4ba-42e3-ae87-b9076630f2b6.zip"},
			want: "e894fa50-c4ba-42e3-ae87-b9076630f2b6",
			wantErr: false,
		},
		{
			name: "get upload uuid from s3 key with .txt.gz",
			args: args{s: "upload/2020/10/7/23/e894fa50-c4ba-42e3-ae87-b9076630f2b6.txt.gz"},
			want: "e894fa50-c4ba-42e3-ae87-b9076630f2b6",
			wantErr: false,
		},
		{
			name: "fail if input is empty",
			args: args{s: ""},
			want: "",
			wantErr: true,
		},
		{
			name: "get upload uuid from s3 key with .txt",
			args: args{s: "upload/2020/10/7/e894fa50-c4ba-42e3-ae87-b9076630f2b6.txt"}, // hour part is missing
			want: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uploadUUID(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("uploadUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("uploadUUID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileType(t *testing.T) {
	type args struct {
		objectKey string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   string
		wantErr bool
	}{
		{name: "passing in a text file", args: args{objectKey: "test.txt"}, want: 450 * 1024, want1: "txt", wantErr: false},
		{name: "passing in a zip file", args: args{objectKey: "test.zip"}, want: 40 * 1024, want1: "zip", wantErr: false},
		{name: "passing in a gzipped file", args: args{objectKey: "test.txt.gz"}, want: 40 * 1024, want1: "gz", wantErr: false},
		{name: "passing something else", args: args{objectKey: "test.jpg"}, want: 0, want1: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := fileType(tt.args.objectKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fileType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("fileType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
