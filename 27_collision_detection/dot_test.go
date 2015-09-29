package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDot(t *testing.T) {
	// assert := assert.New(t)

	var dot Dot = newDot()
	if assert.NotNil(t, dot) {
		assert.Equal(t, 20, dot.width)
	}

}
