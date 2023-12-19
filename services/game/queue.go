package game

import "fmt"

type Queue struct {
	queue []*player
}

type WaitList interface {
	FindAndDelete(p *player) error
	Empty() bool
	Add(p *player) error
	Pop() (*player, error)
	Remove(p *player) error
}

func NewWaitList() WaitList {
	return &Queue{queue: make([]*player, 0)}
}

func (l *Queue) Remove(p *player) error {
	for i := 0; i < len(l.queue); i++ {
		if l.queue[i] == p {
			l.queue = append(l.queue[:i], l.queue[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("player not found")
}

func (l *Queue) Pop() (*player, error) {
	if l.Empty() {
		return nil, fmt.Errorf("empty list")
	}
	p := l.queue[0]
	l.queue = l.queue[1:]
	return p, nil
}

func (l *Queue) Add(p *player) error {
	for i := range l.queue {
		if l.queue[i].GetId() == p.GetId() {
			return fmt.Errorf("already in the waiting list")
		}
	}
	l.queue = append(l.queue, p)
	return nil
}

func (l *Queue) FindAndDelete(p *player) error {
	for i := range l.queue {
		if p.GetId() == l.queue[i].GetId() {
			l.queue = append(l.queue[:i], l.queue[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("player %s not in the list", p.GetId())
}

func (l *Queue) Empty() bool {
	return len(l.queue) == 0
}
