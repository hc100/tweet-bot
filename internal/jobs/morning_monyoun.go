package jobs

import "time"

type MorningMonyounJob struct {
	location *time.Location
}

func NewMorningMonyounJob(location *time.Location) MorningMonyounJob {
	return MorningMonyounJob{location: location}
}

func (j MorningMonyounJob) Name() string {
	return "morning-monyoun"
}

func (j MorningMonyounJob) NextRun(after time.Time) (time.Time, bool) {
	local := after.In(j.location)
	next := time.Date(local.Year(), local.Month(), local.Day(), 6, 0, 0, 0, j.location)
	if !next.After(local) {
		next = next.AddDate(0, 0, 1)
	}

	return next, true
}

func (j MorningMonyounJob) BuildPost(time.Time) (string, error) {
	return "もにゅん", nil
}
