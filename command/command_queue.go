package command

import "sync"

type CommandQueue struct {
	mu       sync.Mutex
	commands []Command
}

func NewCommandQueue() *CommandQueue {
	return &CommandQueue{
		commands: make([]Command, 0),
	}
}

func (q *CommandQueue) Enqueue(cmd Command) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.commands = append(q.commands, cmd)
}

func (q *CommandQueue) Drain() []Command {
	q.mu.Lock()
	defer q.mu.Unlock()

	drained := q.commands
	q.commands = make([]Command, 0)

	return drained
}
