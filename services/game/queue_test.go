package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueAdd(t *testing.T) {
	queue := NewMemoryWaitList()

	p := NewPlayer(nil, nil)
	key := "<10>"
	queue.Add(key, p)

	assert.Equal(t, p, queue.queues[key][0])
}
func TestQueueRemove(t *testing.T) {
	queue := NewMemoryWaitList()

	p := NewPlayer(nil, nil)
	keys := []string{"<10>", "<20>", "<30>"}
	for _, key := range keys {
		queue.Add(key, p)
	}

	queue.Remove(p)

	for _, key := range keys {
		assert.Equal(t, 0, len(queue.queues[key]))
	}
}
func TestQueuePop(t *testing.T) {
	queue := NewMemoryWaitList()

	players := []*player{
		NewPlayer(nil, nil),
		NewPlayer(nil, nil),
		NewPlayer(nil, nil),
	}
	key := "<10>"

	for _, p := range players {
		queue.Add(key, p)
	}

	t.Run("Nil list pop", func(t *testing.T) {
		_, err := queue.Pop("<30>")
		assert.Error(t, err)
	})

	t.Run("Pop order", func(t *testing.T) {
		for _, p := range players {
			pPopped, _ := queue.Pop(key)
			assert.Equal(t, p, pPopped)
		}
	})
	t.Run("Pop on empty list", func(t *testing.T) {
		_, err := queue.Pop("<10>")
		assert.Error(t, err)
	})
}
