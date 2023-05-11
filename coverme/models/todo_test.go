package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinishUnfinish(t *testing.T) {
	todo := Todo{ID: 1, Title: "Test Todo 1", Content: "Some content", Finished: false}
	todo.MarkFinished()
	assert.Equal(t, todo.Finished, true)
	todo.MarkUnfinished()
	assert.Equal(t, todo.Finished, false)
}
