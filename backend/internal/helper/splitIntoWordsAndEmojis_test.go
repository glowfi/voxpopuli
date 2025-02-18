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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := helper.SplitIntoWordsAndEmojis(tt.args.input)
			assert.Equal(t, tt.wantResult, gotResult, "expect results to match")
		})
	}
}
