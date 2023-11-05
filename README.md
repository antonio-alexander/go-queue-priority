# go-queue-priority (github.com/antonio-alexander/go-queue-priority)

go-queue-priority is a FIFO data structure that implements the [go-queue](github.com/antonio-alexander/go-queue) interfaces/functionalit with the addition of being able to prioritize the items being enqueued. This adds to the missing functionality with Go channels: once you've placed data into a channel, there's no way to change the order of the items within without first removing all of the items from the channel.

Similar to the ability to "peek" into a queue, a priority queue will re-arrange (non-destructively) the contents of a queue depending on the priority of the item being enqueued. An item with a higher priority will be dequeued sooner than an item with a lower priority.

Here are some common situations where go-queue-priority functionality would be advantageous:

- If you have a command architecture and although commands are handled serially, you need the ability to send multiple commands to the front of the queue (EnqueueInFront() _can_ do this, but not quite the same way)
- If your consuming data, say historical data and realtime data in the same queue, you can give realtime data a higher priority so that it's consumed BEFORE historical data
- If you want to dynamically sort data using a runtime condition

## Priority Queue Interfaces

go-queue-priority is separated into a high-level/common go module [github.com/antonio-alexander/go-queue-priority](github.com/antonio-alexander/go-queue-priority) where the interfaces (described below) and tests are defined and can be imported/used by anyone attempting to implement those interfaces.

> If it's not obvious, the goal of this separation of ownership of interfaces is used such that anyone using queues depend on the interface, not the implementation

Keep in mind that some of these functions are dependent on the underlying implementation; for example overflow and capacity will have different output depending on if the queue is finite or infinite.

The priority queue implements all of the interfaces of goqueue such as:

- goqueue.Owner
- goqueue.GarbageCollecter
- goqueue.Length
- goqueue.Event
- goqueue.Peeker
- goqueue.Dequeuer
- goqueue.Enqueuer
- finite.EnqueueLossy
- finite.Resizer
- finite.Capacity

Generally, the functionality for all of these are maintained; finite queues will overflow when the queue is full while infinite queues will not. EnqueueLossy too, if the queue is full, it'll push the newest items out with some caveats (see below).

The _biggest_ caveat for the priority queue is the priority queue functionality. The "secret sauce" of the priority queue is that everything that is placed in a queue is wrapped in this data type:

```go
type Wrapper struct {
    Priority   int         `json:"priority"`
    EnqueuedAt int64       `json:"enqueued_at"`
    Item       interface{} `json:"item"`
}
```

The priority queue works by sorting this wrapper, it sorts first by priority, then it sorts by when it was enqueued. The sorting is very anticlimatic:

```go
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
```

Each time data is enqueued into the priority queue, it's sorted, this ensures that whenever a dequeue occurs the correct item is returned (and if multiple items are dequeued, they're dequeued in the right order). The sorting ensures that items that have a greater priority are at the front of the queue and items that have the same priority are enqueued with the oldest items being in the front.

The priority queue provides a _new_ interface, specific for priority queues, this is simply the goqueue.Enqueue interface with the addition of an _optional_ priority. If priority is provided it'll set the priority value in the wrapper.

> The existing goqueue.Enqueue interface simply enqueues items with the default priority of 0, but follows the same rules

```go
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
```

## Patterns

The priority queue doesn't _really_ enable any "new" patterns. All of the old/existing patterns are sill valid and work (e.g., producer/consumer), but in general a priority queue can add more functionality to a producer/consumer.

One of the limitations to producer/consumer is that all items are treated equally, EVEN with EnqueueLossy and EnqueueInFront, they lack intelligence, they can put items in the front or throw away the oldest item to put it a new item, but that's very destructive and equally not so useful.

Priority queues are interesting because you have more flexibility than the EnqueueLossy/EnqueueInFront. Queues are generally immutable; immutable in the sense that once data is in a queue, it's position doesn't really change. Priority queues upend this idea in the sense that the order of the items in the queue changes as you give them different priority.

For example, lets say you are utilizing the producer/consumer pattern, but instead of something where everything is equal (e.g., a historian), your producer/consumer is used to handle commands. Although a command handler doesn't have to serially handle commands, there's a lot of reasons why you may choose to. Let say that it takes "time" to handle your commands, what if you have 10 items queued up, but you need a certain command to execute first because it invalidates the commands that come after it, or you have multiple commands that are simply more important than those other 10.

> It's not a stretch to say that you could implement the above functionality _without_ using a priority queue, but you'd have to do a lot of work: it'd be janky

A priority queue would allow you the ability to constantly sort and change the order of the items inside the queue and allow you the additional functionalit without having to do anything but set the priority. In this implementation of the producer/consumer design pattern, only the producer changes:

```go
const specialPriority int = 10
var queue goqueuepriority.PriorityEnquerer

tProduce := time.NewTicker(time.Second)
defer tProduce.Stop()
tProduceSpecial := time.NewTicker(5*time.Second)
defer tProduceSpecial.Stop()
for {
    select {
    case <-tProduceSpecial.C:
        item := "!!"
        if overflow := queue.PriorityEnqueue(item, specialPriority); !overflow {
            fmt.Printf("enqueued special: %v\n", tNow)
        }
    case <-tProduce.C:
        tNow := time.Now()
        if overflow := queue.PriorityEnqueue(tNow); !overflow {
            fmt.Printf("enqueued: %v\n", tNow)
        }
    }
}
```

The only difference between this producer/consumer and a _normal_ producer consumer, is that even though the bangs are enqueued every five seconds, you should see them as soon as they are enqueued versus waiting for the timestamps to come through.

```go
var queue goqueue.Dequeuer

tConsume := time.NewTicker(time.Second)
defer tConsume.Stop()
for {
    select {
    case <-tConsume.C:
        if item, underflow := queue.Dequeue(); !underflow {
            switch v := item.(type) {
                default:
                    fmt.Printf("unsupported type: %T\n", v)
                case time.Time, *time.Time, string:
                    fmt.Printf("dequeued: %v\n", v)
            }
        }
    }
}
```

Either of the above patterns can be modified (slightly) to be event-based using the _goqueue.Event_ interface and using signals.

## Testing

## Finite Priority Queue
