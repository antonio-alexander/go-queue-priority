package priority

// DefaultPriority is the priority assigned to any items that are
// enqueued that doing have an assigned priority
const DefaultPriority int = 0

// Wrapper is used to provide context to items that are placed into
// the queue, each item that you add to the priority queue is placed
// within this wrapper
type Wrapper struct {
	Priority   int         `json:"priority"`
	EnqueuedAt int64       `json:"enqueued_at"`
	Item       interface{} `json:"item"`
}

// PriorityEnqueuer describes an interface for enqueueing items
// with priority
type PriorityEnqueuer interface {
	//PriorityEnqueue can be used to enqueue a single item
	// with an optional priority; this can be a drop-in replacement
	// for Enqueue()
	PriorityEnqueue(item interface{}, priority ...int) (overflow bool)

	//PriorityEnqueueMultiple can be used to enqueue zero or more items
	// with an optional priority; this can be a drop-in replacement for
	// Enqueue() a single priority can be provided OR a priority for
	// each item can be provided
	PriorityEnqueueMultiple(items []interface{}, priority ...int) (itemsRemaining []interface{}, overflow bool)
}

type ByPriority []*Wrapper

func (b ByPriority) Len() int           { return len(b) }
func (b ByPriority) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByPriority) Less(i, j int) bool { return b[i].Priority > b[j].Priority }

type ByEnqueuedAt []*Wrapper

func (b ByEnqueuedAt) Len() int { return len(b) }
func (b ByEnqueuedAt) Swap(i, j int) {
	if b[i].Priority == b[j].Priority {
		b[i], b[j] = b[j], b[i]
	}
}
func (b ByEnqueuedAt) Less(i, j int) bool { return b[i].EnqueuedAt < b[j].EnqueuedAt }
