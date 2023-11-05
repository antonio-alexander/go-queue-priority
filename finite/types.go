package priorityfinite

type PriorityEnqueueLossy interface {
	PriorityEnqueueLossy(item interface{}, priority ...int) (interface{}, bool)
}
