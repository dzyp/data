package v1

import (
	r "github.com/dzyp/data/trees/rangetree"
)

type orderedList struct {
	nodes     []r.Entry // ordered list of nodes
	dimension int
}

func newOrderedList(dimension int) *orderedList {
	return &orderedList{
		dimension: dimension,
	}
}

func (self *orderedList) insert(entries ...r.Entry) {
	if len(entries) == 0 {
		return
	}

	self.nodes, _ = Entries(self.nodes).Merge(entries...)
}

func (self *orderedList) copy() itree {
	cp := make([]r.Entry, len(self.nodes))
	copy(cp, self.nodes)
	return &orderedList{
		nodes:     cp,
		dimension: self.dimension,
	}
}

func (self *orderedList) query(query r.Query, result *result) {
	bounds := query.GetDimensionalBounds(self.dimension)
	i := Entries(self.nodes).Index(bounds.Low(), self.dimension)

	for i := i; i < len(self.nodes); i++ {
		if self.nodes[i].GetDimensionalValue(self.dimension) < bounds.High() {
			result.AddEntry(self.nodes[i])
		} else {
			break
		}
	}
}

func (self *orderedList) all(result *result) {
	result.AddEntry(self.nodes...)
}

func (self *orderedList) len() int {
	return len(self.nodes)
}

func (self *orderedList) remove(entries ...r.Entry) {
	self.nodes = Entries(self.nodes).Remove(entries...)
}

func (self *orderedList) Insert(entries ...r.Entry) {
	Entries(entries).Sort(1) // this can only happen on the first dimension

	self.insert(entries...)
}
