package v1

import (
	"log"
	"sort"

	r "github.com/dzyp/data/trees/rangetree"
)

type byDimension int

func (dimension byDimension) Sort(entries []r.Entry) {
	es := &entrySorter{
		entries:   entries,
		dimension: int(dimension),
	}

	sort.Sort(es)
}

type entrySorter struct {
	entries   []r.Entry
	dimension int
}

func (self *entrySorter) Len() int {
	return len(self.entries)
}

func (self *entrySorter) Swap(i, j int) {
	self.entries[i], self.entries[j] = self.entries[j], self.entries[i]
}

func (self *entrySorter) Less(i, j int) bool {
	return self.entries[i].Less(self.entries[j], self.dimension)
}

type entriesWrapper struct {
	entries                 []r.Entry
	groups                  map[int][]r.Entry
	sortedDimensionalValues []int
}

func (self *entriesWrapper) getEntriesAtValue(value int) []r.Entry {
	if entries, ok := self.groups[value]; ok {
		return entries
	}

	return nil
}

func (self *entriesWrapper) find(value int) int {
	return sort.SearchInts(self.sortedDimensionalValues, value)
}

func (self *entriesWrapper) getSortedValues() []int {
	return self.sortedDimensionalValues
}

func (self *entriesWrapper) median() int {
	return self.sortedDimensionalValues[len(self.sortedDimensionalValues)/2]
}

func (self *entriesWrapper) isLastValue() bool {
	if self.groups == nil {
		return true
	}

	return len(self.groups) <= 1
}

func (self *entriesWrapper) lastValue() r.Entry {
	if !self.isLastValue() {
		return nil
	}

	key := self.sortedDimensionalValues[0]

	return self.groups[key][0]
}

/*
splits the entity wrapper into the left and right halves
*/
func (self *entriesWrapper) split(index int) (*entriesWrapper, *entriesWrapper) {
	if self.sortedDimensionalValues == nil {
		return nil, nil
	}

	if index == -1 {
		index = len(self.sortedDimensionalValues) / 2
	}

	leftSortedList := self.sortedDimensionalValues[0:index]
	rightSortedList :=
		self.sortedDimensionalValues[index:len(self.sortedDimensionalValues)]

	leftGroup := make(map[int][]r.Entry)
	rightGroup := make(map[int][]r.Entry)

	midPoint := 0

	for _, entry := range leftSortedList {
		entries, ok := self.groups[entry]
		if !ok {
			log.Fatal(`Error in split, entries not matching.`)
		}

		midPoint += len(entries)
		leftGroup[entry] = entries
	}

	for _, entry := range rightSortedList {
		entries, ok := self.groups[entry]
		if !ok {
			log.Fatal(`Error in split, entries not matching.`)
		}

		rightGroup[entry] = entries
	}

	left, right := &entriesWrapper{
		groups:                  leftGroup,
		sortedDimensionalValues: leftSortedList,
	}, &entriesWrapper{
		groups:                  rightGroup,
		sortedDimensionalValues: rightSortedList,
	}

	left.entries = make([]r.Entry, midPoint)
	copy(left.entries, self.entries[0:midPoint])

	right.entries = make([]r.Entry, len(self.entries)-midPoint)
	copy(right.entries, self.entries[midPoint:len(self.entries)])

	return left, right
}

func (self *entriesWrapper) len() int {
	if self.sortedDimensionalValues == nil {
		return 0
	}

	return len(self.sortedDimensionalValues)
}

/*
This should only be called at the top level entry function so we only
sort once.
*/
func newEntries(entries []r.Entry, dimension int, sort bool) *entriesWrapper {
	if len(entries) == 0 {
		return &entriesWrapper{
			entries: entries,
		}
	}

	if sort {
		byDimension(dimension).Sort(entries)
	}

	sortedDimensionalValues := make([]int, len(entries))
	lastSeen := entries[0].GetDimensionalValue(dimension)
	groups := make(map[int][]r.Entry)
	lastIndex := 0
	var sortedIndex int

	for i := 0; i < len(entries); i++ {
		if entries[i].GetDimensionalValue(dimension) == lastSeen {
			continue
		}

		sortedDimensionalValues[sortedIndex] = lastSeen
		groups[lastSeen] = entries[lastIndex:i]

		lastIndex = i
		lastSeen = entries[i].GetDimensionalValue(dimension)
		sortedIndex++
	}

	sortedDimensionalValues[sortedIndex] = lastSeen
	groups[lastSeen] = entries[lastIndex:len(entries)]

	return &entriesWrapper{
		entries:                 entries,
		groups:                  groups,
		sortedDimensionalValues: sortedDimensionalValues[0 : sortedIndex+1],
	}
}
