package jobs

import (
	"fmt"
	"hash/fnv"
	"time"
)

type CountdownTo2027Job struct {
	location *time.Location
	target   time.Time
}

func NewCountdownTo2027Job(location *time.Location) CountdownTo2027Job {
	return CountdownTo2027Job{
		location: location,
		target:   time.Date(2027, 1, 1, 0, 0, 0, 0, location),
	}
}

func (j CountdownTo2027Job) Name() string {
	return "countdown-to-2027"
}

func (j CountdownTo2027Job) NextRun(after time.Time) (time.Time, bool) {
	local := after.In(j.location)
	if !local.Before(j.target) {
		return time.Time{}, false
	}

	candidate := j.scheduledTime(local)
	if candidate.After(local) {
		return candidate, true
	}

	nextDay := local.AddDate(0, 0, 1)
	candidate = j.scheduledTime(nextDay)
	if !candidate.Before(j.target) {
		return time.Time{}, false
	}

	return candidate, true
}

func (j CountdownTo2027Job) BuildPost(now time.Time) (string, error) {
	local := now.In(j.location)
	if !local.Before(j.target) {
		return "", fmt.Errorf("countdown target already reached")
	}

	remaining := j.target.Sub(local)
	days := remaining / (24 * time.Hour)
	remaining -= days * 24 * time.Hour
	hours := remaining / time.Hour
	remaining -= hours * time.Hour
	minutes := remaining / time.Minute
	remaining -= minutes * time.Minute
	seconds := remaining / time.Second

	return fmt.Sprintf(
		"2027年まで%d日%d時間%d分%d秒です。",
		days,
		hours,
		minutes,
		seconds,
	), nil
}

func (j CountdownTo2027Job) scheduledTime(day time.Time) time.Time {
	localDay := day.In(j.location)
	year, month, date := localDay.Date()
	base := time.Date(year, month, date, 10, 0, 0, 0, j.location)

	window := int64((12 * time.Hour) / time.Second)
	offset := deterministicOffsetSeconds(localDay, window)
	return base.Add(time.Duration(offset) * time.Second)
}

func deterministicOffsetSeconds(day time.Time, limit int64) int64 {
	if limit <= 0 {
		return 0
	}

	h := fnv.New64a()
	_, _ = h.Write([]byte(day.Format("2006-01-02")))
	return int64(h.Sum64() % uint64(limit))
}
