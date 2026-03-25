package jobs

import (
	"testing"
	"time"
)

func TestCountdownScheduledTimeIsStableAndWithinWindow(t *testing.T) {
	loc := time.FixedZone("JST", 9*60*60)
	job := NewCountdownTo2027Job(loc)
	day := time.Date(2026, 3, 25, 0, 0, 0, 0, loc)

	first := job.scheduledTime(day)
	second := job.scheduledTime(day)

	if !first.Equal(second) {
		t.Fatalf("scheduledTime should be stable: %v != %v", first, second)
	}
	if first.Hour() < 10 || first.Hour() >= 22 {
		t.Fatalf("scheduled time hour = %d, want 10 <= hour < 22", first.Hour())
	}
}

func TestCountdownBuildPost(t *testing.T) {
	loc := time.FixedZone("JST", 9*60*60)
	job := NewCountdownTo2027Job(loc)
	now := time.Date(2026, 12, 31, 23, 59, 30, 0, loc)

	got, err := job.BuildPost(now)
	if err != nil {
		t.Fatalf("BuildPost returned error: %v", err)
	}

	want := "2027年まで0日0時間0分30秒です。"
	if got != want {
		t.Fatalf("BuildPost = %q, want %q", got, want)
	}
}

func TestCountdownStopsAfterTarget(t *testing.T) {
	loc := time.FixedZone("JST", 9*60*60)
	job := NewCountdownTo2027Job(loc)
	after := time.Date(2027, 1, 1, 0, 0, 0, 0, loc)

	_, ok := job.NextRun(after)
	if ok {
		t.Fatal("NextRun should stop after target")
	}
}
