package loading

import (
	"testing"

	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/stretchr/testify/assert"
)

func TestParallelPlan_IsAnalyzed(t *testing.T) {
	assert.True(t, config.Engine.IsAnalyzed(&ParallelPlan{}))
}

func TestParallelPlan_IsExecutable(t *testing.T) {
	assert.Nil(t, config.Engine.IsExecutable(&ParallelPlan{}))
}
