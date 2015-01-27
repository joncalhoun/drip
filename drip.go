package drip

import (
	"errors"
	"sync"
	"time"
)

type Bucket struct {
	Capacity     int
	DripInterval time.Duration
	PerDrip      int
	consumed     int
	started      bool
	kill         chan bool
	m            sync.Mutex
}

func (b *Bucket) Start() error {
	if b.started {
		return errors.New("Bucket was already started.")
	}

	ticker := time.NewTicker(b.DripInterval)
	b.started = true
	b.kill = make(chan bool, 1)

	go func() {
		for {
			select {
			case <-ticker.C:
				b.m.Lock()
				b.consumed -= b.PerDrip
				if b.consumed < 0 {
					b.consumed = 0
				}
				b.m.Unlock()
			case <-b.kill:
				return
			}
		}
	}()

	return nil
}

func (b *Bucket) Stop() error {
	if !b.started {
		return errors.New("Bucket was never started.")
	}

	b.kill <- true

	return nil
}

func (b *Bucket) Consume(amt int) error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.Capacity-b.consumed < amt {
		return errors.New("Not enough capacity.")
	}
	b.consumed += amt
	return nil
}
