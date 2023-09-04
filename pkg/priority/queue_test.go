package priority_test

import (
	"testing"

	"github.com/jdbann/forestry/pkg/priority"
	"gotest.tools/v3/assert"
)

func TestQueue(t *testing.T) {
	var queue *priority.Queue[string]

	runStep(t, "starts empty", func(t *testing.T) {
		queue = priority.NewQueue[string](5)
		assert.Equal(t, queue.Len(), 0)
	})

	runStep(t, "pushing items in order pops items in order", func(t *testing.T) {
		queue.Push("first", 1)
		queue.Push("second", 2)
		queue.Push("third", 3)
		assert.Equal(t, queue.Pop(), "first")
		assert.Equal(t, queue.Pop(), "second")
		assert.Equal(t, queue.Pop(), "third")
	})

	runStep(t, "pushing items out of order pops items in order", func(t *testing.T) {
		queue.Push("third", 3)
		queue.Push("second", 1)
		queue.Push("first", 2)
		assert.Equal(t, queue.Pop(), "second")
		assert.Equal(t, queue.Pop(), "first")
		assert.Equal(t, queue.Pop(), "third")
	})
}

func runStep(t *testing.T, name string, fn func(t *testing.T)) {
	if !t.Run(name, fn) {
		t.FailNow()
	}
}
