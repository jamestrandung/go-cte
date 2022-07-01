package vat

import (
	"testing"

	"github.com/jamestrandung/go-die/sample/config"
	"github.com/stretchr/testify/assert"
)

func TestComputer_IsRegistered(t *testing.T) {
	assert.True(t, config.Engine.IsRegistered(Amount{}))
}
