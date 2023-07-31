package scheduledjob

import "time"

type SchedulerJob struct {
	ID    string
	Opt   map[string]string
	RunAt time.Time
}
