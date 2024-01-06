package game

import "fmt"

type Queue struct {
	queues map[string][]*player
}

type WaitList interface {
	FindAndDelete(p *player) error
	Empty() bool
	Add(key string, p *player) error
	Pop(key string) (*player, error)
	Remove(p *player) error
}

func NewMemoryWaitList() *Queue {
	return &Queue{queues: make(map[string][]*player, 0)}
}

func (l *Queue) Remove(p *player) error {
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

func (l *Queue) Pop(key string) (*player, error) {
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

func (l *Queue) Add(key string, p *player) error {
	queue, ok := l.queues[key]
	if !ok {
		l.queues[key] = append(l.queues[key], p)
		return nil
	}

	for i := range queue {
		if queue[i].GetId() == p.GetId() {
			return fmt.Errorf("already in the waiting list")
		}
	}
	l.queues[key] = append(l.queues[key], p)
	return nil
}

func (l *Queue) FindAndDelete(p *player) error {
	for key := range l.queues {
		for i := range l.queues[key] {
			if p.GetId() == l.queues[key][i].GetId() {
				l.queues[key] = append(l.queues[key][:i], l.queues[key][i+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("player %s not in the list", p.GetId())
}

func (l *Queue) Empty() bool {
	return len(l.queues) == 0
}
