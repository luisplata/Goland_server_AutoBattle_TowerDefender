package game

import "time"

type GameClock struct {
	tickDuration time.Duration
	lastTickTime time.Time
}

func NewGameClock(tickDurationMs int) *GameClock {
	return &GameClock{
		tickDuration: time.Duration(tickDurationMs) * time.Millisecond,
		lastTickTime: time.Now().UTC(),
	}
}

func (c *GameClock) ShouldTick() bool {
	now := time.Now().UTC()
	if now.Sub(c.lastTickTime) >= c.tickDuration {
		c.lastTickTime = now
		return true
	}
	return false
}
