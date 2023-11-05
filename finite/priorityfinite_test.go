package priorityfinite_test

import (
	"testing"
	"time"

	goqueue "github.com/antonio-alexander/go-queue"
	goqueuepriority "github.com/antonio-alexander/go-queue-priority"
	goqueuepriorityfinite "github.com/antonio-alexander/go-queue-priority/finite"
	finite "github.com/antonio-alexander/go-queue/finite"

	goqueuepriorityfinite_tests "github.com/antonio-alexander/go-queue-priority/finite/tests"
	finite_tests "github.com/antonio-alexander/go-queue/finite/tests"
	goqueue_tests "github.com/antonio-alexander/go-queue/tests"
)

const (
	mustTimeout time.Duration = time.Second
	mustRate    time.Duration = time.Millisecond
)

func TestFiniteQueue(t *testing.T) {
	t.Run("Test Enqueue", finite_tests.TestEnqueue(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Enqueue Multiple", finite_tests.TestEnqueueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Enqueue Event", finite_tests.TestEnqueueEvent(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Event
	} {
		return goqueuepriorityfinite.New(size)
	}))

	t.Run("Test Resize", finite_tests.TestResize(t, func(size int) interface {
		finite.Capacity
		goqueue.Enqueuer
		goqueue.Owner
		finite.Resizer
	} {
		return goqueuepriorityfinite.New(size)
	}))
	//REVIEW: how to fix this functionality?
	// t.Run("Test Enqueue Lossy", finite_tests.TestEnqueueLossy(t, func(size int) interface {
	// 	goqueue.Owner
	// 	finite.EnqueueLossy
	// } {
	// 	return goqueuepriorityfinite.New(size)
	// }))
	t.Run("Test Capacity", finite_tests.TestCapacity(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		finite.Capacity
	} {
		return goqueuepriorityfinite.New(size)
	}))
}

func TestQueue(t *testing.T) {
	t.Run("Test Dequeue", goqueue_tests.TestDequeue(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Dequeue Event", goqueue_tests.TestDequeueEvent(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Dequeuer
		goqueue.Enqueuer
		goqueue.Event
		goqueue.Owner
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Dequeue Multiple", goqueue_tests.TestDequeueMultiple(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Flush", goqueue_tests.TestFlush(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Peek", goqueue_tests.TestPeek(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Peeker
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Peek From Head", goqueue_tests.TestPeekFromHead(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Peeker
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Length", goqueue_tests.TestLength(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Garbage Collect", goqueue_tests.TestGarbageCollect(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.GarbageCollecter
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return goqueuepriorityfinite.New(size)
	}))
	//
	t.Run("Test Queue", goqueue_tests.TestQueue(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Asynchronous", goqueue_tests.TestAsync(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return goqueuepriorityfinite.New(size)
	}))
}

func TestPriorityFiniteQueue(t *testing.T) {
	t.Run("Test Priority Enqueue", goqueuepriorityfinite_tests.TestPriorityEnqueue(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Dequeuer
		goqueuepriority.PriorityEnqueuer
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Priority Enqueue Event", goqueuepriorityfinite_tests.TestPriorityEnqueueEvent(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Dequeuer
		goqueue.Event
		goqueuepriority.PriorityEnqueuer
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Priority Enqueue Lossy", goqueuepriorityfinite_tests.TestPriorityEnqueueLossy(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Dequeuer
		goqueuepriority.PriorityEnqueuer
		goqueuepriorityfinite.PriorityEnqueueLossy
	} {
		return goqueuepriorityfinite.New(size)
	}))
	t.Run("Test Priority Enqueue Lossy Event", goqueuepriorityfinite_tests.TestPriorityEnqueueLossyEvent(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Dequeuer
		goqueue.Event
		goqueuepriority.PriorityEnqueuer
		goqueuepriorityfinite.PriorityEnqueueLossy
	} {
		return goqueuepriorityfinite.New(size)
	}))
}
