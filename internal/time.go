package fleetlock

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// days = ${weekdays//\'/\"}
// start_time = "$start_time"
// length_minutes = $length

type timeLock struct {
	active    bool
	days      []time.Weekday
	startTime time.Time
	length    time.Duration
}

func NewTimeLock() (*timeLock, error) {
	lock := &timeLock{
		active: false,
		days:   []time.Weekday{},
		length: time.Duration(2) * time.Hour,
	}

	tl, err := strconv.ParseBool(os.Getenv("TIMELOCK"))
	if err != nil {
		lock.active = false
	}
	lock.active = tl

	daysEnc, ok := os.LookupEnv("TIMELOCK_DAYS")
	if ok {
		var days []string
		err := json.Unmarshal([]byte(daysEnc), &days)
		if err != nil {
			return nil, err
		}

		for _, d := range days {
			switch strings.ToLower(d) {
			case "mon":
				lock.days = append(lock.days, time.Monday)
			case "tue":
				lock.days = append(lock.days, time.Tuesday)
			case "wed":
				lock.days = append(lock.days, time.Wednesday)
			case "thu":
				lock.days = append(lock.days, time.Thursday)
			case "fri":
				lock.days = append(lock.days, time.Friday)
			case "sat":
				lock.days = append(lock.days, time.Saturday)
			case "sun":
				lock.days = append(lock.days, time.Sunday)
			default:
				return nil, fmt.Errorf("unsupported day, cannot parse: %s", d)
			}
		}
	} else {
		lock.days = append(lock.days, time.Saturday)
	}

	du, ok := os.LookupEnv("TIMELOCK_DURATION")
	if ok {
		dui, err := strconv.Atoi(du)
		if err != nil {
			return nil, err
		}
		lock.length = time.Duration(dui) * time.Minute
	}

	sts, ok := os.LookupEnv("TIMELOCK_STARTTIME")
	layout := "15:04"
	if !ok {
		sts = "21:00"
	}
	st, err := time.Parse(layout, sts)
	if err != nil {
		return nil, err
	}
	lock.startTime = st
	return lock, nil
}

func (l *timeLock) IsActive() bool {
	return l.active
}

func (l *timeLock) CanLock(n time.Time) bool {
	// Set day to the given time to be able to compare the timeframe
	// Remove one month and day since the day and month are 1 by default
	st := l.startTime.AddDate(n.Year(), int(n.Month())-1 , n.Day()-1)
	et := st.Add(l.length)

	day := n.Weekday()
	// We need to account for a condition when the timeslot reaches into the next day
	if n.Before(st) || n.After(et) {
		// Return if time frame from the previous day is before the given time
		if n.After(et.AddDate(0,0,-1)) {
			return false
		}
		// Get week day of previous day
		day = n.AddDate(0,0,-1).Weekday()
	}

	return slices.Contains(l.days, day)
}