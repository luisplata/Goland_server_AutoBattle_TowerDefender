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

// TicksPerSecond returns how many ticks occur per second based on the clock's tick duration.
func (c *GameClock) TicksPerSecond() int {
	if c.tickDuration <= 0 {
		return 1
	}
	return int(time.Second / c.tickDuration)
}
