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
	groups                  map[r.Entry][]r.Entry
	sortedDimensionalValues []r.Entry
}

func (self *entriesWrapper) getEntriesAtValue(value r.Entry) []r.Entry {
	if entries, ok := self.groups[value]; ok {
		return entries
	}

	return nil
}

func (self *entriesWrapper) getSortedValues() []r.Entry {
	return self.sortedDimensionalValues
}

func (self *entriesWrapper) median() r.Entry {
	if self.sortedDimensionalValues == nil {
		return nil
	}

	return self.sortedDimensionalValues[len(self.sortedDimensionalValues)/2]
}

func (self *entriesWrapper) isLastValue() bool {
	if self.groups == nil {
		return true
	}

	return len(self.groups) <= 1
}

/*
splits the entity wrapper into the left and right halves
*/
func (self *entriesWrapper) split() (*entriesWrapper, *entriesWrapper) {
	if self.sortedDimensionalValues == nil {
		return nil, nil
	}

	median := len(self.sortedDimensionalValues) / 2
	leftSortedList := self.sortedDimensionalValues[0:median]
	rightSortedList :=
		self.sortedDimensionalValues[median:len(self.sortedDimensionalValues)]

	leftGroup := make(map[r.Entry][]r.Entry)
	rightGroup := make(map[r.Entry][]r.Entry)

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

	return &entriesWrapper{
			entries:                 self.entries[0:midPoint],
			groups:                  leftGroup,
			sortedDimensionalValues: leftSortedList,
		}, &entriesWrapper{
			entries:                 self.entries[midPoint:len(self.entries)],
			groups:                  rightGroup,
			sortedDimensionalValues: rightSortedList,
		}
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

	sortedDimensionalValues := make([]r.Entry, len(entries))
	lastSeen := entries[0]
	groups := make(map[r.Entry][]r.Entry)
	lastIndex := 0
	var sortedIndex int

	for i := 0; i < len(entries); i++ {
		if entries[i].EqualAtDimension(lastSeen, dimension) {
			continue
		}
		sortedDimensionalValues[sortedIndex] = lastSeen
		groups[lastSeen] = entries[lastIndex:i]

		lastIndex = i
		lastSeen = entries[i]
		sortedIndex++
	}

	sortedDimensionalValues[sortedIndex] = lastSeen
	groups[lastSeen] = entries[lastIndex:len(entries)]

	log.Printf(`dimensional values: %+v`, sortedDimensionalValues)

	return &entriesWrapper{
		entries:                 entries,
		groups:                  groups,
		sortedDimensionalValues: sortedDimensionalValues[0 : sortedIndex+1],
	}
}
