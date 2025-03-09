package helper_test

import (
	"testing"

	"github.com/glowfi/voxpopuli/backend/internal/helper"
	"github.com/stretchr/testify/assert"
)

func Test_SplitStringIntoCustomEmojisAndWords(t *testing.T) {
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
		{
			name: "case 8:",
			args: args{
				input: "Science &amp; Technology 🧬:hey: :hello:",
			},
			wantResult: []string{"Science &amp; Technology 🧬", ":hey:", " ", ":hello:"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := helper.SplitStringIntoCustomEmojisAndWords(tt.args.input)
			assert.Equal(t, tt.wantResult, gotResult, "expect results to match")
		})
	}
}

func Test_SplitStringIntoStandardEmojisAndWords(t *testing.T) {
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
				input: "🧬Science &amp; Technology 🧬",
			},
			wantResult: []string{"🧬", "Science &amp; Technology ", "🧬"},
		},
		{
			name: "case 2:",
			args: args{
				input: "🧬🧬 hello🧬hey",
			},
			wantResult: []string{"🧬", "🧬", " hello", "🧬", "hey"},
		},
		{
			name: "case 3:",
			args: args{
				input: " hello hey    hi🧬  🧬",
			},
			wantResult: []string{" hello hey    hi", "🧬", "  ", "🧬"},
		},
		{
			name: "case 4:",
			args: args{
				input: "🙂🙂  hey🙂🙂 hello",
			},
			wantResult: []string{"🙂", "🙂", "  hey", "🙂", "🙂", " hello"},
		},
		{
			name: "case 5:",
			args: args{
				input: "  hey🙂🙂 hello🙂",
			},
			wantResult: []string{"  hey", "🙂", "🙂", " hello", "🙂"},
		},
		{
			name: "case 6:",
			args: args{
				input: "JEFFREE⭐",
			},
			wantResult: []string{"JEFFREE", "⭐"},
		},
		{
			name: "case 7:",
			args: args{
				input: "#️⃣1️⃣ WORLD FIRST!",
			},
			wantResult: []string{
				"#️⃣",
				"1️⃣",
				" WORLD FIRST!",
			},
		},
		{
			name: "case 8:",
			args: args{
				input: "ぁé́́hello1️⃣1️⃣world 1️⃣hello  1️⃣hey 1️⃣",
			},
			wantResult: []string{"ぁé́́hello", "1️⃣", "1️⃣", "world ", "1️⃣", "hello  ", "1️⃣", "hey ", "1️⃣"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := helper.SplitStringIntoStandardEmojisAndWords(tt.args.input)
			assert.Equal(t, tt.wantResult, gotResult, "expect results to match")
		})
	}
}
