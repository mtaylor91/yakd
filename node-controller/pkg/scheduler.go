package pkg

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Schedule struct {
	executeDone       bool
	executeFunc       func()
	executeNext       *Schedule
	executeTime       time.Time
	executeRepeatedly time.Duration
}

type Scheduler struct {
	schedules   chan<- *Schedule
	executeNext *Schedule
}

func New() *Scheduler {
	schedules := make(chan *Schedule)

	scheduler := &Scheduler{schedules, nil}
	go scheduler.run(schedules)

	return scheduler
}

func (scheduler *Scheduler) add(schedule *Schedule) {
	if scheduler.executeNext == nil {
		scheduler.executeNext = schedule
	} else if schedule.executeTime.Before(scheduler.executeNext.executeTime) {
		schedule.executeNext = scheduler.executeNext
		scheduler.executeNext = schedule
	} else {
		executeNext := scheduler.executeNext

		for executeNext.executeNext != nil &&
			executeNext.executeNext.executeTime.Before(schedule.executeTime) {
			executeNext = executeNext.executeNext
		}

		schedule.executeNext = executeNext.executeNext
		executeNext.executeNext = schedule
	}
}

func (scheduler *Scheduler) executeNextSchedule() {
	if scheduler.executeNext.executeTime.Before(time.Now()) {
		schedule := scheduler.executeNext
		scheduler.executeNext = schedule.executeNext
		schedule.executeFunc()
		if schedule.executeRepeatedly != 0 {
			schedule.executeTime.Add(schedule.executeRepeatedly)
			scheduler.add(schedule)
		} else {
			schedule.executeDone = true
		}
	}
}

func (scheduler *Scheduler) run(schedules <-chan *Schedule) {
	log.Info("Scheduler started")
	for {
		scheduler.schedule(schedules)
	}
}

func (scheduler *Scheduler) schedule(schedules <-chan *Schedule) {
	start := time.Now()
	sleep := 1 * time.Hour

	if scheduler.executeNext != nil {
		if start.Before(scheduler.executeNext.executeTime) {
			sleep = scheduler.executeNext.executeTime.Sub(start)
		} else {
			sleep = 0
		}
	}

	select {
	case schedule := <-schedules:
		scheduler.add(schedule)
		log.Info("Schedule added")
	case <-time.After(sleep):
	}

	if scheduler.executeNext != nil {
		scheduler.executeNextSchedule()
	}

	finish := time.Now()
	elapsed := finish.Sub(start)
	log.WithFields(log.Fields{
		"elapsed": elapsed,
	}).Info("Scheduling next")
}

func (s *Scheduler) ScheduleOnce(f func(), t time.Time) {
	s.schedules <- &Schedule{
		executeDone:       false,
		executeFunc:       f,
		executeNext:       nil,
		executeTime:       t,
		executeRepeatedly: 0,
	}
}

func (s *Scheduler) ScheduleRepeat(f func(), t time.Time, r time.Duration) {
	s.schedules <- &Schedule{
		executeDone:       false,
		executeFunc:       f,
		executeNext:       nil,
		executeTime:       t,
		executeRepeatedly: r,
	}
}
