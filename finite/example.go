package priorityfinite

import (
	goqueue "github.com/antonio-alexander/go-queue"
)

//TODO: add ExamplePriorityEnqueue

func ExamplePriorityEnqueueLossy(queue PriorityEnqueueLossy, value *goqueue.Example, priorities ...int) (*goqueue.Example, bool) {
	item, discarded := queue.PriorityEnqueueLossy(value, priorities...)
	if !discarded {
		return nil, false
	}
	return goqueue.ExampleConvertSingle(item), true
}
