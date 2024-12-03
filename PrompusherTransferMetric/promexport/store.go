package promexport

import (
	"container/list"
	"sync"
)

type PromMetricList struct {
	sync.RWMutex

	maxSize int
	li      *list.List
}

func NewPromMetricList(maxSize int) *PromMetricList {
	return &PromMetricList{
		maxSize: maxSize,
		li:      list.New(),
	}
}

// default PushFront
func (pml *PromMetricList) PushFrontBatch(vs []interface{}) bool {
	if pml.li.Len() >= pml.maxSize {
		return false
	}
	pml.Lock()
	for _, item := range vs {
		pml.li.PushFront(item)
	}
	pml.Unlock()
	return true
}

func (pml *PromMetricList) PushOne(item interface{}) bool {
	if pml.li.Len() >= pml.maxSize {
		return false
	}
	pml.Lock()
	pml.li.PushFront(item)
	pml.Unlock()

	return true
}

func (pml *PromMetricList) PopBackBatch(max int) []interface{} {

	pml.Lock()
	count := pml.li.Len()
	if count == 0 {
		pml.Unlock()
		return nil
	}
	if count > max {
		count = max
	}
	items := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		item := pml.li.Remove(pml.li.Back())
		items = append(items, item)
	}
	pml.Unlock()

	return items
}

func (pml *PromMetricList) PopBackAll() []interface{} {

	pml.Lock()
	count := pml.li.Len()
	if count == 0 {
		pml.Unlock()
		return nil
	}

	items := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		item := pml.li.Remove(pml.li.Back())
		items = append(items, item)
	}
	pml.Unlock()

	return items
}

func (pml *PromMetricList) Len() int {
	pml.RLock()
	defer pml.RUnlock()

	return pml.li.Len()
}
