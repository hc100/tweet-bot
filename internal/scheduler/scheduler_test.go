package scheduler

import (
	"context"
	"testing"
	"time"
)

type stubPoster struct {
	posts []string
	err   error
}

func (p *stubPoster) Post(_ context.Context, text string) error {
	if p.err != nil {
		return p.err
	}
	p.posts = append(p.posts, text)
	return nil
}

type stubJob struct {
	name string
	run  time.Time
	text string
	ok   bool
}

func (j stubJob) Name() string { return j.name }

func (j stubJob) NextRun(time.Time) (time.Time, bool) {
	return j.run, j.ok
}

func (j stubJob) BuildPost(time.Time) (string, error) {
	return j.text, nil
}

func TestNextJobSelectsEarliestRun(t *testing.T) {
	now := time.Date(2026, 3, 25, 9, 0, 0, 0, time.UTC)
	s := New(&stubPoster{}, time.UTC, []Job{
		stubJob{name: "late", run: now.Add(2 * time.Hour), ok: true},
		stubJob{name: "early", run: now.Add(1 * time.Hour), ok: true},
	})

	runAt, job, ok := s.nextJob(now)
	if !ok {
		t.Fatal("nextJob returned ok=false")
	}
	if job.Name() != "early" {
		t.Fatalf("job = %s, want early", job.Name())
	}
	if !runAt.Equal(now.Add(1 * time.Hour)) {
		t.Fatalf("runAt = %v, want %v", runAt, now.Add(1*time.Hour))
	}
}

func TestRunReturnsErrorWhenNoJobsAreRunnable(t *testing.T) {
	s := New(&stubPoster{}, time.UTC, nil)
	s.now = func() time.Time { return time.Date(2026, 3, 25, 9, 0, 0, 0, time.UTC) }

	err := s.Run(context.Background())
	if err == nil {
		t.Fatal("Run returned nil error")
	}
	if err.Error() != "no runnable jobs configured" {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
}
