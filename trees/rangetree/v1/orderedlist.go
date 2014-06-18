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

	nodesAdded := 0
	cpEntries := make([]r.Entry, len(entries))
	copy(cpEntries, entries)

	newNodes := make([]r.Entry, 0, len(entries)+len(self.nodes))

	var newNode r.Entry
	var oldNode r.Entry

	for {
		if len(self.nodes) == 0 && len(cpEntries) == 0 {
			break
		} else if len(self.nodes) == 0 {
			newNodes = append(newNodes, cpEntries...)
			nodesAdded += len(cpEntries)
			break
		} else if len(cpEntries) == 0 {
			newNodes = append(newNodes, self.nodes...)
			break
		}

		newNode = cpEntries[0]
		oldNode = self.nodes[0]

		if newNode.Less(oldNode, 1) {
			newNodes = append(newNodes, newNode)
			cpEntries = Entries(cpEntries).RemoveAt(0)
			nodesAdded++
		} else if oldNode.Less(newNode, 1) {
			newNodes = append(newNodes, oldNode)
			self.nodes = Entries(self.nodes).RemoveAt(0)
			nodesAdded++
		} else { //equal
			//println(`EQUAL`)
			newNodes = append(newNodes, newNode) // we override the old value
			self.nodes = Entries(self.nodes).RemoveAt(0)
			cpEntries = Entries(cpEntries).RemoveAt(0)
			// we don't add to nodes added here as this number of rebalancing purposes
		}
	}

	self.nodes = newNodes

	return nodesAdded
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

func (self *orderedList) Insert(entries ...r.Entry) {
	Entries(entries).Sort(1) // this can only happen on the first dimension

	self.insert(entries...)
}
