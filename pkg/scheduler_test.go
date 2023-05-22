package pkg

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAddSchedules(t *testing.T) {
	s1 := &Schedule{executeTime: time.Now().Add(1 * time.Second)}
	s2 := &Schedule{executeTime: time.Now().Add(2 * time.Second)}
	s3 := &Schedule{executeTime: time.Now().Add(3 * time.Second)}
	s4 := &Schedule{executeTime: time.Now().Add(4 * time.Second)}
	scheduler := &Scheduler{}
	scheduler.add(s4)
	scheduler.add(s2)
	scheduler.add(s3)
	scheduler.add(s1)
	assert.True(t, scheduler.executeNext == s1)
	assert.True(t, scheduler.executeNext.executeNext == s2)
	assert.True(t, scheduler.executeNext.executeNext.executeNext == s3)
	assert.True(t, scheduler.executeNext.executeNext.executeNext.executeNext == s4)
}

func TestScheduler(t *testing.T) {
	// Create a new scheduler
	scheduler := New()

	// Make channels to receive results
	c1 := make(chan struct{})
	c2 := make(chan struct{})
	c3 := make(chan struct{})

	// Schedule tasks to submit results
	scheduler.ScheduleOnce(func() {
		log.Info("Execute 3")
		close(c3)
	}, time.Now().Add(3*time.Millisecond))
	scheduler.ScheduleOnce(func() {
		log.Info("Execute 2")
		close(c2)
	}, time.Now().Add(2*time.Millisecond))
	scheduler.ScheduleOnce(func() {
		log.Info("Execute 1")
		close(c1)
	}, time.Now().Add(1*time.Millisecond))

	// Wait for results from channels
	<-c1
	<-c2
	<-c3
}
