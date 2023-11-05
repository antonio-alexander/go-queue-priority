package internal

import (
	"time"

	goqueuepriority "github.com/antonio-alexander/go-queue-priority"
)

// RotateLeft can be used to perform an in-place rotation
// left of a slice of empty interface
func RotateLeft(items []*goqueuepriority.Wrapper) {
	if len(items) > 1 {
		copy(items, append(items[1:], items[:1]...))
	}
}

// RotateRight can be used to perform an in-place rotation
// right of a slice of empty interface
func RotateRight(items []*goqueuepriority.Wrapper) {
	if len(items) < 1 {
		return
	}
	copy(items, append(items[len(items)-1:], items[:len(items)-1]...))
}

// Enqueue can be used  to add an item to the back of a queue while maintaining
// it's capacity (e.g. in-place) it will return true if the queue is full
func Enqueue(items []*goqueuepriority.Wrapper, item *goqueuepriority.Wrapper) ([]*goqueuepriority.Wrapper, bool) {
	if len(items) >= cap(items) {
		return items, true
	}
	return append(items, item), false
}

// Dequeue can be used to remove an item from the queue and reduce its
// capacity by one
func Dequeue(items []*goqueuepriority.Wrapper) (interface{}, []*goqueuepriority.Wrapper, bool) {
	if len(items) <= 0 {
		return nil, items, true
	}
	item := items[0]
	items[0] = nil
	RotateLeft(items)
	if len(items) > 0 {
		items = items[:len(items)-1] //truncate the slice
	}
	return item.Item, items, false
}

// DequeueMultiple will return a number of items less than or equal to the value of
// n while maintaining the input data on the second slice of interface, it will return
// true if there are no items to dequeue
func DequeueMultiple(n int, wrappers []*goqueuepriority.Wrapper) ([]interface{}, []*goqueuepriority.Wrapper, bool) {
	var length int

	//get the length of the data, underflow if no data, then
	// check to see if n is negative or greater than -1
	if length = len(wrappers); length <= 0 {
		return nil, wrappers, true
	}
	if n > length {
		n = length
	}
	items := make([]interface{}, 0, n)
	//TODO: do this at once instead of using dequeue
	// over and over again
	for i := 0; i < n; i++ {
		var item interface{}

		item, wrappers, _ = Dequeue(wrappers)
		items = append(items, item)
	}
	return items, wrappers, false
}

// SendSignal will perform a non-blocking send with or without
// a timeout depending on whether ConfigSignalTimeout is greater
// than 0
func SendSignal(signal chan struct{}, timeout ...time.Duration) bool {
	if len(timeout) > 0 {
		select {
		case <-time.After(timeout[0]):
		case signal <- struct{}{}:
			return true
		}
		return false
	}
	select {
	default:
	case signal <- struct{}{}:
		return true
	}
	return false
}
