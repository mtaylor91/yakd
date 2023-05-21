package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddSchedules(t *testing.T) {
	s1 := &Schedule{executeTime: time.Now().Add(1 * time.Second)}
	s2 := &Schedule{executeTime: time.Now().Add(2 * time.Second)}
	s3 := &Schedule{executeTime: time.Now().Add(3 * time.Second)}
	s4 := &Schedule{executeTime: time.Now().Add(4 * time.Second)}
	s := New()
	s.add(s4)
	s.add(s2)
	s.add(s3)
	s.add(s1)
	assert.True(t, s.executeNext == s1)
	assert.True(t, s.executeNext.executeNext == s2)
	assert.True(t, s.executeNext.executeNext.executeNext == s3)
	assert.True(t, s.executeNext.executeNext.executeNext.executeNext == s4)
}
