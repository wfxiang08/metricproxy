package forwarder

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/signalfuse/signalfxproxy/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestForwarderBasicBufferedForwarder(t *testing.T) {
	upto := uint32(10)
	f := newBasicBufferedForwarder(100, upto, "aname", 1)
	assert.Equal(t, "aname", f.Name(), "Mismatched name")
	for i := uint32(0); i < upto+uint32(1); i++ {
		f.DatapointsChannel() <- nil
	}
	points := f.blockingDrainUpTo()
	assert.Equal(t, upto, uint32(len(points)), "Expect a full set")
	points = f.blockingDrainUpTo()
	assert.Equal(t, 1, len(points), "Expect one point")
}

func TestForwarderStopForwarder(t *testing.T) {
	f := newBasicBufferedForwarder(100, uint32(10), "aname", 1)
	f.DatapointsChannel() <- nil
	f.DatapointsChannel() <- nil
	log.Info("Hello")
	seenPointsChan := make(chan int, 2)
	// nil should make it stop itself
	f.start(func(dp []core.Datapoint) error {
		defer func() {
			seenPointsChan <- len(dp)
		}()
		return errors.New("unable to process")
	})
	seenPoints := <-seenPointsChan
	f.stop()
	assert.Equal(t, 2, seenPoints, "Expect two points")
	assert.NotEqual(t, nil, f.start(nil), "Shouldn't be able to start twice")
	f.blockingDrainStopChan <- true
	f.stop() // Should not block if chan has an item
}
