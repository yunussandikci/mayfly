package common

import (
	"slices"
	"time"

	"github.com/elliotchance/pie/v2"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler interface {
	CreateOrUpdateTask(tag string, date time.Time, task func() error) error
	DeleteTask(tag string) error
}

type scheduler struct {
	config    *Config
	scheduler gocron.Scheduler
}

func NewScheduler(config *Config) Scheduler {
	schedulerInstance := &scheduler{
		config: config,
	}

	cronScheduler, newSchedulerErr := gocron.NewScheduler()
	if newSchedulerErr != nil {
		panic(newSchedulerErr)
	}

	if _, newJobErr := cronScheduler.NewJob(gocron.DurationJob(config.MonitoringInterval), gocron.NewTask(func() {
		mayflyTotalJobs.Set(float64(len(cronScheduler.Jobs())) - 1)
	}), gocron.WithTags("monitoring")); newJobErr != nil {
		panic(newJobErr)
	}

	schedulerInstance.scheduler = cronScheduler
	cronScheduler.Start()

	return schedulerInstance
}

func (s *scheduler) CreateOrUpdateTask(tag string, date time.Time, task func() error) error {
	job := pie.Of(s.scheduler.Jobs()).Filter(func(job gocron.Job) bool {
		return slices.Contains(job.Tags(), tag)
	}).First()

	if job != nil {
		_, updateErr := s.scheduler.Update(job.ID(), gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(date)), gocron.NewTask(task), gocron.WithTags(tag))

		return updateErr
	}

	_, jobErr := s.scheduler.NewJob(gocron.OneTimeJob(
		gocron.OneTimeJobStartDateTime(date)), gocron.NewTask(task), gocron.WithTags(tag))

	return jobErr
}

func (s *scheduler) DeleteTask(tag string) error {
	job := pie.Of(s.scheduler.Jobs()).Filter(func(job gocron.Job) bool {
		return slices.Contains(job.Tags(), tag)
	}).First()

	if job != nil {
		return s.scheduler.RemoveJob(job.ID())
	}

	return nil
}
