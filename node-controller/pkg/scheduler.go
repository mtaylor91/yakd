package pkg

import (
	"time"
)

type Schedule struct {
	executeDone bool
	executeFunc func()
	executeNext *Schedule
	executeTime time.Time
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

func (scheduler *Scheduler) run(schedules <-chan *Schedule) {
	for {
		scheduler.schedule(schedules)
	}
}

func (scheduler *Scheduler) schedule(schedules <-chan *Schedule) {
	schedule := <-schedules
	scheduler.add(schedule)

	finish := time.Now()
	if scheduler.executeNext != nil &&
		scheduler.executeNext.executeTime.Before(finish) {
		schedule := scheduler.executeNext
		scheduler.executeNext = schedule.executeNext
		schedule.executeFunc()
		schedule.executeDone = true
	}
}

func (s *Scheduler) ScheduleOnce(f func(), t time.Time) {
	s.schedules <- &Schedule{
		executeDone: false,
		executeFunc: f,
		executeNext: nil,
		executeTime: t,
	}
}
