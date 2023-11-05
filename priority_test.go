package priority_test

import (
	"sort"
	"testing"
	"time"

	goqueuepriority "github.com/antonio-alexander/go-queue-priority"

	"github.com/stretchr/testify/assert"
)

const casef string = "case: %s"

// TestSortByPriorityAndEnqueuedAt is meant to confirm the functionality
// of the sorting algorithms that should sort by priority initially, and
// then (without separating the priority) by when it was enqueued
func TestSortByPriorityAndEnqueuedAt(t *testing.T) {
	tNow := time.Now()
	priority := goqueuepriority.DefaultPriority
	cases := map[string]struct {
		iWrappers []*goqueuepriority.Wrapper
		oWrappers []*goqueuepriority.Wrapper
	}{
		"same_priority": {
			iWrappers: []*goqueuepriority.Wrapper{
				{EnqueuedAt: tNow.Add(3 * time.Minute).UnixNano()},
				{EnqueuedAt: tNow.Add(2 * time.Minute).UnixNano()},
				{EnqueuedAt: tNow.Add(1 * time.Minute).UnixNano()},
			},
			oWrappers: []*goqueuepriority.Wrapper{
				{EnqueuedAt: tNow.Add(1 * time.Minute).UnixNano()},
				{EnqueuedAt: tNow.Add(2 * time.Minute).UnixNano()},
				{EnqueuedAt: tNow.Add(3 * time.Minute).UnixNano()},
			},
		},
		"different_priority_same_time": {
			iWrappers: []*goqueuepriority.Wrapper{
				{
					Priority:   priority,
					EnqueuedAt: tNow.UnixNano(),
				},
				{
					Priority:   priority + 2,
					EnqueuedAt: tNow.UnixNano(),
				},
				{
					Priority:   priority + 3,
					EnqueuedAt: tNow.UnixNano(),
				},
			},
			oWrappers: []*goqueuepriority.Wrapper{
				{
					Priority:   priority + 3,
					EnqueuedAt: tNow.UnixNano(),
				},
				{
					Priority:   priority + 2,
					EnqueuedAt: tNow.UnixNano(),
				},
				{
					Priority:   priority,
					EnqueuedAt: tNow.UnixNano(),
				},
			},
		},
		"different_priority_different_times": {
			iWrappers: []*goqueuepriority.Wrapper{
				{
					Priority:   priority,
					EnqueuedAt: tNow.Add(4 * time.Second).UnixNano(),
				},
				{
					Priority:   priority + 1,
					EnqueuedAt: tNow.Add(3 * time.Second).UnixNano(),
				},
				{
					Priority:   priority + 1,
					EnqueuedAt: tNow.Add(2 * time.Second).UnixNano(),
				},
				{
					Priority:   priority + 2,
					EnqueuedAt: tNow.Add(0 * time.Second).UnixNano(),
				},
			},
			oWrappers: []*goqueuepriority.Wrapper{
				{
					Priority:   priority + 2,
					EnqueuedAt: tNow.Add(0 * time.Second).UnixNano(),
				},
				{
					Priority:   priority + 1,
					EnqueuedAt: tNow.Add(2 * time.Second).UnixNano(),
				},
				{
					Priority:   priority + 1,
					EnqueuedAt: tNow.Add(3 * time.Second).UnixNano(),
				},
				{
					Priority:   priority,
					EnqueuedAt: tNow.Add(4 * time.Second).UnixNano(),
				},
			},
		},
	}
	for cDesc, c := range cases {
		wrappers := make([]*goqueuepriority.Wrapper, len(c.iWrappers))
		copy(wrappers, c.iWrappers)
		sort.Sort(goqueuepriority.ByPriority(wrappers))
		sort.Sort(goqueuepriority.ByEnqueuedAt(wrappers))
		assert.Equal(t, c.oWrappers, wrappers, casef, cDesc)
	}
}
