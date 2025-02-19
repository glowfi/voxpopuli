package helper_test

import (
	"testing"

	"github.com/glowfi/voxpopuli/backend/internal/helper"
	"github.com/stretchr/testify/assert"
)

func Test_splitIntoWordsAndEmojis(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []string
	}{
		{
			name: "case 1:",
			args: args{
				input: "hello     :nin::yang:hello  helloworld",
			},
			wantResult: []string{"hello     ", ":nin:", ":yang:", "hello  helloworld"},
		},
		{
			name: "case 2:",
			args: args{
				input: "hello     :nin::yang:hello  helloworld:nin:",
			},
			wantResult: []string{"hello     ", ":nin:", ":yang:", "hello  helloworld", ":nin:"},
		},
		{
			name: "case 3:",
			args: args{
				input: "\"hello     \":nin::yang:hello  helloworld:nin:",
			},
			wantResult: []string{"\"hello     \"", ":nin:", ":yang:", "hello  helloworld", ":nin:"},
		},
		{
			name: "case 4:",
			args: args{
				input: "hello :ni::n::ja:",
			},
			wantResult: []string{"hello ", ":ni:", ":n:", ":ja:"},
		},
		{
			name: "case 5:",
			args: args{
				input: "desc1 :e1::e1: :ce1::ce1:",
			},
			wantResult: []string{"desc1 ", ":e1:", ":e1:", " ", ":ce1:", ":ce1:"},
		},
		{
			name: "case 6:",
			args: args{
				input: "desc2 :e2::e2: :ce2::ce2: desc2",
			},
			wantResult: []string{"desc2 ", ":e2:", ":e2:", " ", ":ce2:", ":ce2:", " desc2"},
		},
		{
			name: "case 7:",
			args: args{
				input: ":e3::ce3:hello",
			},
			wantResult: []string{":e3:", ":ce3:", "hello"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := helper.SplitIntoWordsAndEmojis(tt.args.input)
			assert.Equal(t, tt.wantResult, gotResult, "expect results to match")
		})
	}
}
