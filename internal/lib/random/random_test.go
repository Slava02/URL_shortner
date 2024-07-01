package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	type argsTestLen struct {
		len int
	}

	testLen := []struct {
		name string
		args argsTestLen
		want int
	}{
		{
			name: "len == 0",
			args: argsTestLen{0},
			want: 0,
		},
		{
			name: "len == 4",
			args: argsTestLen{4},
			want: 4,
		},
	}

	type argsTesEq struct {
		amount, len int
	}

	testEq := []struct {
		name string
		args argsTesEq
	}{
		{
			name: "testEq len == 100",
			args: argsTesEq{100, 4},
		},
	}

	for _, tt := range testLen {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomString(tt.args.len)
			str2 := NewRandomString(tt.args.len)

			if tt.args.len == 0 {
				assert.Equal(t, str1, "")
				assert.Equal(t, str2, "")
			} else {
				assert.Len(t, str1, tt.args.len)
				assert.Len(t, str2, tt.args.len)
				assert.NotEqual(t, str1, str2)
			}
		})
	}

	for _, tt := range testEq {
		t.Run(tt.name, func(t *testing.T) {
			var randStr []string
			for i := 0; i < tt.args.len; i++ {
				randStr = append(randStr, NewRandomString(tt.args.len))
			}
			for i := 0; i < len(randStr)-1; i++ {
				assert.NotEqual(t, randStr[i], randStr[i+1])
			}
		})
	}
}
