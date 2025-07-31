package test

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestIntegration_Setup(t *testing.T) {
	a := 2
	b := 4
	assert.Equal(t, a+b, 6)
}
