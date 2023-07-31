package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/lusis/go-rundeck/pkg/rundeck"
	scheduledjob "github.com/machtwatch/catalyst-go-skeleton/domain/scheduler"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/tracer"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/machtwatch/catalystdk/go/trace"
)

const (
	segmentRundeckRunJob  = tracer.SegmentRundeck + "RunJob"
	segmentRundeckKillJob = tracer.SegmentRundeck + "KillJob"
)

type (
	Scheduler struct {
		client      *rundeck.Client
		projectName string
	}

	SchedulerConfig struct {
		RundeckUrl      string
		RundeckAPIToken string
		RundeckProject  string
	}

	ISchedulerExecution interface {
		RunJob(ctx context.Context, job scheduledjob.SchedulerJob) (int, error)
		KillJob(ctx context.Context, id int) error
	}
)

func NewRundeckConnection(config *SchedulerConfig) ISchedulerExecution {
	return CreateScheduler(config)
}

// CreateScheduler initialize connection to Rundeck and return a Scheduler
func CreateScheduler(config *SchedulerConfig) *Scheduler {
	client, err := rundeck.NewTokenAuthClient(config.RundeckAPIToken, config.RundeckUrl)
	if err != nil {
		log.Fatalf("rundeck.NewTokenAuthClient() error on connecting to rundeck: %v", err)
	}

	user, err := client.GetCurrentUserProfile()
	if err != nil {
		log.Fatalf("client.GetCurrentUserProfile() - error on connecting to rundeck: %v", err)
	}

	log.Infof("successfully connected to rundeck with profile: %v", user.Login)
	return &Scheduler{
		client:      client,
		projectName: config.RundeckProject,
	}
}

// RunJob calls Rundeck API to execute a job with specified time and options and returns an execution ID
func (s Scheduler) RunJob(ctx context.Context, job scheduledjob.SchedulerJob) (int, error) {
	ctxRundeck, span := trace.StartSpanFromContext(ctx, segmentRundeckRunJob)
	defer span.End()

	exec, err := s.client.RunJob(job.ID, rundeck.RunJobOpts(job.Opt), rundeck.RunJobRunAt(job.RunAt))
	if err != nil {
		log.StdDebug(ctxRundeck, job, err, "s.client.RunJob(job.ID, rundeck.RunJobOpts(job.Opt), rundeck.RunJobRunAt(job.RunAt))")
		return 0, fmt.Errorf("error running scheduled job, error: %v", err)
	}

	log.Infof("job ID %v successfully scheduled at %v with execution ID: %v. Payload: %+v", job.ID, job.RunAt, exec.ID, job.Opt)
	return exec.ID, nil
}

// KillJob calls Rundeck API to kill a scheduled job and returns an error message if any
func (s Scheduler) KillJob(ctx context.Context, id int) error {
	ctxRundeck, span := trace.StartSpanFromContext(ctx, segmentRundeckKillJob)
	defer span.End()

	_, err := s.client.AbortExecution(id, func(m *map[string]string) error { return nil })
	if err != nil {
		log.StdDebug(ctxRundeck, id, err, "s.client.AbortExecution(id, func(m *map[string]string) error { return nil })")
		return fmt.Errorf("error killing scheduled job, error: %v", err)
	}

	log.Infof("execution ID %v successfully abort at %v", id, time.Now())
	return err
}
