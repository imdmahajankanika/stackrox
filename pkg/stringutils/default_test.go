package stringutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	for _, testCase := range []struct {
		s        string
		defaul   string
		expected string
	}{
		{"blah", "default", "blah"},
		{"blah", "", "blah"},
		{"", "default", "default"},
		{"", "", ""},
	} {
		c := testCase
		t.Run(fmt.Sprintf("%+v", c), func(t *testing.T) {
			assert.Equal(t, c.expected, OrDefault(c.s, c.defaul))
		})
	}
}

func TestPointerDefault(t *testing.T) {
	blah := "blah"
	empty := ""
	for _, testCase := range []struct {
		s        *string
		defaul   string
		expected string
	}{
		{&blah, "default", "blah"},
		{&blah, "", "blah"},
		{&empty, "default", "default"},
		{&empty, "", ""},
		{nil, "", ""},
		{nil, "default", "default"},
	} {
		c := testCase
		t.Run(fmt.Sprintf("%+v", c), func(t *testing.T) {
			assert.Equal(t, c.expected, PointerOrDefault(c.s, c.defaul))
		})
	}
}
