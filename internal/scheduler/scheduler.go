package scheduler

import (
	"context"
	"errors"
	"log"
	"time"
)

type Poster interface {
	Post(ctx context.Context, text string) error
}

type Job interface {
	Name() string
	NextRun(after time.Time) (time.Time, bool)
	BuildPost(now time.Time) (string, error)
}

type Scheduler struct {
	poster   Poster
	location *time.Location
	jobs     []Job
	now      func() time.Time
	timer    func(time.Duration) (<-chan time.Time, func() bool)
}

func New(poster Poster, location *time.Location, jobs []Job) *Scheduler {
	return &Scheduler{
		poster:   poster,
		location: location,
		jobs:     jobs,
		now:      time.Now,
		timer:    newTimer,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	for {
		now := s.now().In(s.location)
		runAt, job, ok := s.nextJob(now)
		if !ok {
			return errors.New("no runnable jobs configured")
		}

		wait := time.Until(runAt)
		if wait < 0 {
			wait = 0
		}
		log.Printf("next job=%s run_at=%s", job.Name(), runAt.Format(time.RFC3339))

		timerC, stop := s.timer(wait)
		select {
		case <-ctx.Done():
			stop()
			return ctx.Err()
		case <-timerC:
		}

		executedAt := s.now().In(s.location)
		text, err := job.BuildPost(executedAt)
		if err != nil {
			log.Printf("job=%s build post failed: %v", job.Name(), err)
			continue
		}

		if err := s.poster.Post(ctx, text); err != nil {
			log.Printf("job=%s post failed: %v", job.Name(), err)
			continue
		}

		log.Printf("job=%s posted text=%q", job.Name(), text)
	}
}

func (s *Scheduler) nextJob(now time.Time) (time.Time, Job, bool) {
	var selected Job
	var selectedTime time.Time

	for _, job := range s.jobs {
		runAt, ok := job.NextRun(now)
		if !ok {
			continue
		}
		if selected == nil || runAt.Before(selectedTime) {
			selected = job
			selectedTime = runAt
		}
	}

	if selected == nil {
		return time.Time{}, nil, false
	}

	return selectedTime, selected, true
}

func newTimer(d time.Duration) (<-chan time.Time, func() bool) {
	t := time.NewTimer(d)
	return t.C, t.Stop
}
