package priorityfinite

import (
	"sort"
	"sync"
	"time"

	internal "github.com/antonio-alexander/go-queue-priority/internal"

	goqueue "github.com/antonio-alexander/go-queue"
	priorityqueue "github.com/antonio-alexander/go-queue-priority"
	finite "github.com/antonio-alexander/go-queue/finite"
)

type queueFinite struct {
	sync.RWMutex
	signalIn  chan struct{}
	signalOut chan struct{}
	data      []*priorityqueue.Wrapper
}

func New(size int) interface {
	goqueue.Owner
	goqueue.GarbageCollecter
	goqueue.Length
	goqueue.Event
	goqueue.Peeker
	goqueue.Dequeuer
	goqueue.Enqueuer
	finite.EnqueueLossy
	finite.Resizer
	finite.Capacity
	priorityqueue.PriorityEnqueuer
	PriorityEnqueueLossy
} {
	if size < 1 {
		size = 1
	}
	return &queueFinite{
		signalIn:  make(chan struct{}, size),
		signalOut: make(chan struct{}, size),
		data:      make([]*priorityqueue.Wrapper, 0, size),
	}
}

func (q *queueFinite) enqueueLossy(items []*priorityqueue.Wrapper, itemToEnqueue *priorityqueue.Wrapper) (interface{}, bool) {
	//KIM: this works off of the idea that the items slice is already
	// sorted
	lowestPriority := items[0].Priority
	if itemToEnqueue.Priority < lowestPriority {
		return nil, true
	}
	itemDiscarded := items[0]
	items[0] = itemToEnqueue
	sort.Sort(priorityqueue.ByPriority(items))
	sort.Sort(priorityqueue.ByEnqueuedAt(items))
	return itemDiscarded.Item, false
}

func (q *queueFinite) Close() []interface{} {
	q.Lock()
	defer q.Unlock()

	remainingElements, _, _ := internal.DequeueMultiple(len(q.data), q.data)
	if q.signalIn != nil {
		select {
		default:
			close(q.signalIn)
		case <-q.signalIn:
		}
	}
	if q.signalOut != nil {
		select {
		default:
			close(q.signalOut)
		case <-q.signalOut:
		}
	}
	q.data, q.signalIn, q.signalOut = nil, nil, nil
	return remainingElements
}

func (q *queueFinite) GarbageCollect() {
	q.Lock()
	defer q.Unlock()

	//create a new slice to hold the data copy the data
	// from the old slice to the new slice and set the
	// internal data to be the new slice
	data := make([]*priorityqueue.Wrapper, 0, cap(q.data))
	copy(data, q.data)
	q.data = data
}

func (q *queueFinite) Resize(newSize int) []interface{} {
	q.Lock()
	defer q.Unlock()

	var discardedItems []interface{}

	//ensure that no operations occur if the size hasn't changed,
	// if there's a need to remove items, remove them, then copy the old
	// data to the newly created slice, create new signal channels
	if newSize == cap(q.data) {
		return nil
	}
	if newSize < 1 {
		newSize = 1
	}
	if len(q.data) > newSize {
		discardedItems, q.data, _ = internal.DequeueMultiple(len(q.data)-newSize, q.data)
	}
	data := make([]*priorityqueue.Wrapper, len(q.data), newSize)
	copy(data, q.data[:len(q.data)])
	if q.signalIn != nil {
		select {
		default:
			close(q.signalIn)
		case <-q.signalIn:
		}
	}
	if q.signalOut != nil {
		select {
		default:
			close(q.signalOut)
		case <-q.signalOut:
		}
	}
	q.data = data
	q.signalIn = make(chan struct{}, newSize)
	q.signalOut = make(chan struct{}, newSize)
	return discardedItems
}

func (q *queueFinite) GetSignalIn() <-chan struct{} {
	q.RLock()
	defer q.RUnlock()

	return q.signalIn
}

func (q *queueFinite) GetSignalOut() <-chan struct{} {
	q.RLock()
	defer q.RUnlock()

	return q.signalOut
}

func (q *queueFinite) Dequeue() (interface{}, bool) {
	q.Lock()
	defer q.Unlock()

	var item interface{}
	var underflow bool

	item, q.data, underflow = internal.Dequeue(q.data)
	if underflow {
		return nil, underflow
	}
	internal.SendSignal(q.signalOut)
	return item, false
}

func (q *queueFinite) DequeueMultiple(n int) []interface{} {
	q.Lock()
	defer q.Unlock()

	var items []interface{}
	var underflow bool

	items, q.data, underflow = internal.DequeueMultiple(n, q.data)
	if underflow {
		return nil
	}
	internal.SendSignal(q.signalOut)
	return items
}

func (q *queueFinite) Flush() []interface{} {
	q.Lock()
	defer q.Unlock()

	var items []interface{}
	var underflow bool

	items, q.data, underflow = internal.DequeueMultiple(cap(q.data), q.data)
	if underflow {
		return nil
	}
	internal.SendSignal(q.signalOut)
	return items
}

func (q *queueFinite) Enqueue(item interface{}) bool {
	return q.PriorityEnqueue(item)
}

func (q *queueFinite) EnqueueMultiple(items []interface{}) ([]interface{}, bool) {
	return q.PriorityEnqueueMultiple(items)
}

func (q *queueFinite) EnqueueLossy(item interface{}) (interface{}, bool) {
	return q.PriorityEnqueueLossy(item)
}

func (q *queueFinite) PriorityEnqueue(item interface{}, priorities ...int) bool {
	q.Lock()
	defer q.Unlock()

	var overflow bool

	priority := priorityqueue.DefaultPriority
	if len(priorities) > 0 {
		priority = priorities[0]
	}
	if q.data, overflow = internal.Enqueue(q.data, &priorityqueue.Wrapper{
		Item:       item,
		Priority:   priority,
		EnqueuedAt: time.Now().UnixNano(),
	}); overflow {
		return true
	}
	sort.Sort(priorityqueue.ByPriority(q.data))
	sort.Sort(priorityqueue.ByEnqueuedAt(q.data))
	internal.SendSignal(q.signalIn)
	return false
}

func (q *queueFinite) PriorityEnqueueMultiple(items []interface{}, priorities ...int) ([]interface{}, bool) {
	q.Lock()
	defer q.Unlock()

	var itemEnqueued, overflow bool

	defer func() {
		if itemEnqueued {
			sort.Sort(priorityqueue.ByPriority(q.data))
			sort.Sort(priorityqueue.ByEnqueuedAt(q.data))
			internal.SendSignal(q.signalIn)
		}
	}()
	if len(priorities) != len(items) {
		priority := priorityqueue.DefaultPriority
		if len(priorities) > 0 {
			priority = priorities[0]
		}
		priorities = make([]int, 0, len(items))
		for range items {
			priorities = append(priorities, priority)
		}
	}
	for i, item := range items {
		q.data, overflow = internal.Enqueue(q.data, &priorityqueue.Wrapper{
			Item:       item,
			Priority:   priorities[i],
			EnqueuedAt: time.Now().UnixNano(),
		})
		if overflow {
			return items[i:], overflow
		}
	}
	return nil, false
}

func (q *queueFinite) PriorityEnqueueLossy(item interface{}, priorities ...int) (interface{}, bool) {
	q.Lock()
	defer q.Unlock()

	var overflow bool

	priority := priorityqueue.DefaultPriority
	if len(priorities) > 0 {
		priority = priorities[0]
	}
	wrappedItem := &priorityqueue.Wrapper{
		Item:       item,
		Priority:   priority,
		EnqueuedAt: time.Now().UnixNano(),
	}
	if q.data, overflow = internal.Enqueue(q.data, wrappedItem); !overflow {
		return nil, false
	}
	return q.enqueueLossy(q.data, wrappedItem)
}

func (q *queueFinite) Length() (size int) {
	q.RLock()
	defer q.RUnlock()

	return len(q.data)
}

func (q *queueFinite) Capacity() (capacity int) {
	q.RLock()
	defer q.RUnlock()

	return cap(q.data)
}

func (q *queueFinite) Peek() []interface{} {
	q.RLock()
	defer q.RUnlock()

	items := make([]interface{}, 0, len(q.data))
	for i := 0; i < len(q.data); i++ {
		items = append(items, q.data[i].Item)
	}
	return items
}

func (q *queueFinite) PeekHead() (item interface{}, underflow bool) {
	q.RLock()
	defer q.RUnlock()

	if len(q.data) <= 0 {
		return nil, true
	}
	return q.data[0].Item, false
}

func (q *queueFinite) PeekFromHead(n int) []interface{} {
	q.RLock()
	defer q.RUnlock()

	if len(q.data) == 0 {
		return nil
	}
	if n > len(q.data) {
		n = len(q.data)
	}
	items := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		items = append(items, q.data[i].Item)
	}
	return items
}
