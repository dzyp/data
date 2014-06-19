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

func (self *orderedList) insert(entries ...r.Entry) int {
	if len(entries) == 0 {
		return 0
	}

	var result int
	self.nodes, result = Entries(self.nodes).Merge(entries...)

	return result
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

func (self *orderedList) Insert(entries ...r.Entry) {
	Entries(entries).Sort(1) // this can only happen on the first dimension

	self.insert(entries...)
}
