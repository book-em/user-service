//go:build integration

package test

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestStup(t *testing.T) {
	a := 2
	b := 4
	assert.Equal(t, a+b, 4)
}
