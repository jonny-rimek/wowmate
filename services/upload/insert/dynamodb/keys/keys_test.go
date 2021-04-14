package main

import "testing"

func Test_durationAsPercent(t *testing.T) {
	type args struct {
		dungeonIntimeDuration  float64
		durationInMilliseconds float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "is 100",
			args: args{
				dungeonIntimeDuration:  180,
				durationInMilliseconds: 180,
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := durationAsPercent(tt.args.dungeonIntimeDuration, tt.args.durationInMilliseconds); got != tt.want {
				t.Errorf("durationAsPercent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timed(t *testing.T) {
	type args struct {
		durationInMilliseconds float64
		intimeDuration         float64
		twoChestDuration       float64
		threeChestDuration     float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "three chest",
			args: args{
				durationInMilliseconds: 10,
				intimeDuration:         80,
				twoChestDuration:       60,
				threeChestDuration:     20,
			},
			want: 3,
		},
		{
			name: "two chest",
			args: args{
				durationInMilliseconds: 30,
				intimeDuration:         80,
				twoChestDuration:       60,
				threeChestDuration:     20,
			},
			want: 2,
		},
		{
			name: "intime",
			args: args{
				durationInMilliseconds: 70,
				intimeDuration:         80,
				twoChestDuration:       60,
				threeChestDuration:     20,
			},
			want: 1,
		},
		{
			name: "deplete",
			args: args{
				durationInMilliseconds: 100,
				intimeDuration:         80,
				twoChestDuration:       60,
				threeChestDuration:     20,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timed(tt.args.durationInMilliseconds, tt.args.intimeDuration, tt.args.twoChestDuration, tt.args.threeChestDuration); got != tt.want {
				t.Errorf("timed() = %v, want %v", got, tt.want)
			}
		})
	}
}
