package priority

import (
	"time"

	goqueue "github.com/antonio-alexander/go-queue"
)

func MustPriorityEnqueue(queue PriorityEnqueuer, item interface{}, priority int, done <-chan struct{}, rate time.Duration) bool {
	if overflow := queue.PriorityEnqueue(item, priority); !overflow {
		return overflow
	}
	tEnqueue := time.NewTicker(rate)
	defer tEnqueue.Stop()
	if done != nil {
		for {
			select {
			case <-done:
				return queue.PriorityEnqueue(item, priority)
			case <-tEnqueue.C:
				if overflow := queue.PriorityEnqueue(item, priority); !overflow {
					return overflow
				}
			}
		}
	}
	for {
		<-tEnqueue.C
		if overflow := queue.PriorityEnqueue(item, priority); !overflow {
			return overflow
		}
	}
}

func MustPriorityEnqueueEvent(queue interface {
	PriorityEnqueuer
	goqueue.Event
}, item interface{}, priority int, done <-chan struct{}) bool {
	if overflow := queue.PriorityEnqueue(item, priority); !overflow {
		return false
	}
	signalOut := queue.GetSignalOut()
	if done != nil {
		for {
			select {
			case <-done:
				return queue.PriorityEnqueue(item, priority)
			case <-signalOut:
				if overflow := queue.PriorityEnqueue(item, priority); !overflow {
					return overflow
				}
			}
		}
	}
	for {
		<-signalOut
		if overflow := queue.PriorityEnqueue(item, priority); !overflow {
			return overflow
		}
	}
}

func MustPriorityEnqueueMultiple(queue PriorityEnqueuer, items []interface{}, priorities []int, done <-chan struct{}, rate time.Duration) ([]interface{}, bool) {
	itemsRemaining, overflow := queue.PriorityEnqueueMultiple(items, priorities...)
	if !overflow {
		return nil, false
	}
	items = itemsRemaining
	tEnqueueMultiple := time.NewTicker(rate)
	defer tEnqueueMultiple.Stop()
	if done != nil {
		for {
			select {
			case <-done:
				return queue.PriorityEnqueueMultiple(items, priorities...)
			case <-tEnqueueMultiple.C:
				itemsRemaining, overflow := queue.PriorityEnqueueMultiple(items, priorities...)
				if !overflow {
					return nil, false
				}
				items = itemsRemaining
			}
		}
	}
	for {
		<-tEnqueueMultiple.C
		itemsRemaining, overflow := queue.PriorityEnqueueMultiple(items, priorities...)
		if !overflow {
			return nil, false
		}
		items = itemsRemaining
	}
}

// MustPriorityEnqueueMultipleEvent will attempt to enqueue one or more items, upon initial
// failure, it'll use the event channels/signals to attempt to enqueue items
// KIM: this function doesn't preserve the unit of work and may not be consistent
// with concurent usage (although it is safe)
func MustPriorityEnqueueMultipleEvent(queue interface {
	PriorityEnqueuer
	goqueue.Event
}, items []interface{}, priorities []int, done <-chan struct{}) ([]interface{}, bool) {
	itemsRemaining, overflow := queue.PriorityEnqueueMultiple(items, priorities...)
	if !overflow {
		return nil, false
	}
	items = itemsRemaining
	signalOut := queue.GetSignalOut()
	if done != nil {
		for {
			select {
			case <-done:
				return queue.PriorityEnqueueMultiple(items, priorities...)
			case <-signalOut:
				itemsRemaining, overflow := queue.PriorityEnqueueMultiple(items, priorities...)
				if !overflow {
					return nil, false
				}
				items = itemsRemaining
			}
		}
	}
	for {
		<-signalOut
		itemsRemaining, overflow := queue.PriorityEnqueueMultiple(items, priorities...)
		if !overflow {
			return nil, false
		}
		items = itemsRemaining
	}
}
