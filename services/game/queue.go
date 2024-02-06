package game

import "fmt"

type Queue struct {
	queues map[string][]Client
}

type WaitList interface {
	FindAndDelete(p Client) error
	Empty() bool
	Add(key string, p Client) error
	Pop(key string) (Client, error)
	Remove(p Client) error
}

func NewMemoryWaitList() *Queue {
	return &Queue{queues: make(map[string][]Client, 0)}
}

func (l *Queue) Remove(p Client) error {
	found := false
	for i := range l.queues {
		for j := 0; j < len(l.queues[i]); j++ {
			if l.queues[i][j] == p {
				l.queues[i] = append(l.queues[i][:j], l.queues[i][j+1:]...)
				found = true
			}
		}
	}

	if !found {
		return fmt.Errorf("player not found")
	}
	return nil
}

func (l *Queue) Pop(key string) (Client, error) {
	queue, ok := l.queues[key]
	if !ok {
		return nil, fmt.Errorf("empty list")
	}
	if len(queue) == 0 {
		return nil, fmt.Errorf("empty list")
	}
	p := queue[0]
	l.queues[key] = l.queues[key][1:]
	return p, nil
}

func (l *Queue) Add(key string, p Client) error {
	queue, ok := l.queues[key]
	if !ok {
		l.queues[key] = append(l.queues[key], p)
		return nil
	}

	for i := range queue {
		if queue[i] == p {
			return fmt.Errorf("already in the waiting list")
		}
	}
	l.queues[key] = append(l.queues[key], p)
	return nil
}

func (l *Queue) FindAndDelete(p Client) error {
	for key := range l.queues {
		for i := range l.queues[key] {
			if p == l.queues[key][i] {
				l.queues[key] = append(l.queues[key][:i], l.queues[key][i+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("player %v not in the list", p)
}

func (l *Queue) Empty() bool {
	return len(l.queues) == 0
}
