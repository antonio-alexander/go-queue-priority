package priorityfinite_tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	goqueue "github.com/antonio-alexander/go-queue"
	goqueuepriority "github.com/antonio-alexander/go-queue-priority"
	goqueuepriorityfinite "github.com/antonio-alexander/go-queue-priority/finite"

	"github.com/stretchr/testify/assert"
)

const casef string = "case: %s"

func TestPriorityEnqueue(t *testing.T, rate, timeout time.Duration, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Dequeuer
	goqueuepriority.PriorityEnqueuer
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize       int
			iPriorities []int
			iItems      []interface{}
			oItems      []interface{}
		}{
			"same_priority": {
				iSize:       3,
				iPriorities: []int{1, 1, 1},
				iItems: []interface{}{
					goqueue.Example{Int: 1},
					goqueue.Example{Int: 2},
					goqueue.Example{Int: 3},
				},
				oItems: []interface{}{
					goqueue.Example{Int: 1},
					goqueue.Example{Int: 2},
					goqueue.Example{Int: 3},
				},
			},
			"different_priority": {
				iSize:       3,
				iPriorities: []int{1, 2, 3},
				iItems: []interface{}{
					goqueue.Example{Int: 1},
					goqueue.Example{Int: 2},
					goqueue.Example{Int: 3},
				},
				oItems: []interface{}{
					goqueue.Example{Int: 3},
					goqueue.Example{Int: 2},
					goqueue.Example{Int: 1},
				},
			},
		}
		for cDesc, c := range cases {
			//create queue
			q := newQueue(c.iSize)

			//flush queue
			_ = q.Flush()

			// TODO: enqueue items
			for i, item := range c.iItems {
				ctx, cancel := context.WithTimeout(context.TODO(), timeout)
				defer cancel()
				overflow := goqueuepriority.MustPriorityEnqueue(q, item, c.iPriorities[i],
					ctx.Done(), timeout)
				assert.False(t, overflow, casef, cDesc)
			}

			// TODO: flush items and validate
			ctx, cancel := context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			items := goqueue.MustFlush(q, ctx.Done(), rate)
			assert.Equal(t, c.oItems, items, casef, cDesc)

			//close queue
			q.Close()
		}
		// for cDesc, c := range cases {
		// 	//TODO: enqueue multiple items
		// 	//TODO: dequeue multiple
		// }
		// for cDesc, c := range cases {
		// 	//TODO: enqueue multiple items
		// 	//TODO: flush items and validate
		// }
	}
}

func TestPriorityEnqueueEvent(t *testing.T, rate, timeout time.Duration, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Dequeuer
	goqueue.Event
	goqueuepriority.PriorityEnqueuer
}) func(*testing.T) {
	return func(t *testing.T) {
		//generate size and examples
		size := int(100 * rand.Float64())
		// examples := goqueue.ExampleGenFloat64(size)

		//create queue
		q := newQueue(size)
		defer q.Close()

		//
	}
}

func TestPriorityEnqueueLossy(t *testing.T, rate, timeout time.Duration, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Dequeuer
	goqueuepriority.PriorityEnqueuer
	goqueuepriorityfinite.PriorityEnqueueLossy
}) func(*testing.T) {
	return func(t *testing.T) {
		//generate size and examples
		size := int(100 * rand.Float64())
		// examples := goqueue.ExampleGenFloat64(size)

		//create queue
		q := newQueue(size)
		defer q.Close()

		//
	}
}

func TestPriorityEnqueueLossyEvent(t *testing.T, rate, timeout time.Duration, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Dequeuer
	goqueue.Event
	goqueuepriority.PriorityEnqueuer
	goqueuepriorityfinite.PriorityEnqueueLossy
}) func(*testing.T) {
	return func(t *testing.T) {
		//generate size and examples
		size := int(100 * rand.Float64())
		// examples := goqueue.ExampleGenFloat64(size)

		//create queue
		q := newQueue(size)
		defer q.Close()

		//
	}
}
