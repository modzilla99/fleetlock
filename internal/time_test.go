package fleetlock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCanLock(t *testing.T) {
	layout := "15:04:05 02.01.2006"
	st, err := time.Parse(layout, "20:00:00 01.01.0000")
	if err != nil {
		t.Error(err)
	}

	l := &timeLock{
		active: true,
		days: []time.Weekday{
			time.Monday,
		},
		// Time slots starts at 8pm
		startTime: st,
		// Make time slot end at 2am
		length: time.Duration(6 * time.Hour),
	}

	// Monday 8pm
	t1, _ := time.Parse(layout, "20:00:01 09.01.2006")
	assert.True(t, l.CanLock(t1))

	// Tuesday 1:59am
	t2, _ := time.Parse(layout, "1:59:59 03.01.2006")
	assert.True(t, l.CanLock(t2))

	// Tuesday 2am
	t3, _ := time.Parse(layout, "2:00:01 17.01.2006")
	assert.False(t, l.CanLock(t3))

	// Monday 7:59pm
	t4, _ := time.Parse(layout, "19:59:45 21.01.2006")
	assert.False(t, l.CanLock(t4))
}