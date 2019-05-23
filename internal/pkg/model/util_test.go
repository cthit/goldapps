package model

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestSanitizeEmail(t *testing.T) {
	assert.Equal(t, SanitizeEmail("123abc"), "123abc")
	assert.Equal(t, SanitizeEmail("123aBc"), "123abc")
	assert.Equal(t, SanitizeEmail("123 abc"), "123-abc")
	assert.Equal(t, SanitizeEmail("123-abc"), "123-abc")
	assert.Equal(t, SanitizeEmail("123*abc"), "123abc")
	assert.Equal(t, SanitizeEmail("123öabc"), "123oabc")
}
