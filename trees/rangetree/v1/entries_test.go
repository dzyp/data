package v1

import (
	"testing"

	r "github.com/dzyp/data/trees/rangetree"
)

func checkSortedList(t *testing.T, entries []r.Entry, expected ...*coordinate) {
	if len(entries) != len(expected) {
		t.Errorf(`Expected len: %d, received: %d`, len(expected), len(entries))
		return // don't want to panic
	}

	for i, entry := range entries {
		p := entry.(*point)
		coord := expected[i]

		if p.coordinates[0] != coord.x {
			t.Errorf(`Expected x: %d, received: %d`, coord.x, p.coordinates[0])
		}

		if p.coordinates[1] != coord.y {
			t.Errorf(`Expected y: %d, received: %d`, coord.y, p.coordinates[1])
		}
	}
}

func TestSortOnFirstDimension(t *testing.T) {
	entries := []r.Entry{
		newPoint(2, 0),
		newPoint(4, 0),
		newPoint(3, 0),
		newPoint(1, 0),
	}

	byDimension(1).Sort(entries)

	checkEntries(
		t, entries,
		newCoordinate(1, 0),
		newCoordinate(2, 0),
		newCoordinate(3, 0),
		newCoordinate(4, 0),
	)
}

func TestSortOnSecondDimension(t *testing.T) {
	entries := []r.Entry{
		newPoint(0, 3),
		newPoint(0, 1),
		newPoint(0, 4),
		newPoint(0, 2),
	}

	byDimension(2).Sort(entries)

	checkEntries(
		t, entries,
		newCoordinate(0, 1),
		newCoordinate(0, 2),
		newCoordinate(0, 3),
		newCoordinate(0, 4),
	)
}

func TestEntriesWrapperFirstDimension(t *testing.T) {
	entries := []r.Entry{
		newPoint(0, 1),
		newPoint(1, 0),
		newPoint(1, 1),
		newPoint(0, 0),
	}

	ew := newEntries(entries, 1, true)

	sorted := ew.getSortedValues()

	checkSortedList(t, sorted, newCoordinate(0, 0), newCoordinate(1, 0))

	entries = ew.getEntriesAtValue(sorted[0])
	checkEntries(t, entries, newCoordinate(0, 0), newCoordinate(0, 1))

	entries = ew.getEntriesAtValue(sorted[1])
	checkEntries(t, entries, newCoordinate(1, 0), newCoordinate(1, 1))
}

func TestEntriesWrapperOneValue(t *testing.T) {
	entries := []r.Entry{
		newPoint(0, 1),
	}

	ew := newEntries(entries, 1, true)

	sorted := ew.getSortedValues()
	checkSortedList(t, sorted, newCoordinate(0, 1))

	entries = ew.getEntriesAtValue(sorted[0])
	checkEntries(t, entries, newCoordinate(0, 1))
}

func TestWrapperSplit(t *testing.T) {
	entries := []r.Entry{
		newPoint(0, 1),
		newPoint(1, 1),
		newPoint(0, 0),
		newPoint(1, 0),
	}

	ew := newEntries(entries, 1, true)
	entry := ew.median()

	checkEntries(t, []r.Entry{entry}, newCoordinate(1, 0))

	left, right := ew.split()

	if len(left.groups) != 1 {
		t.Errorf(`Expected len: %d, received: %d`, 1, len(left.groups))
	}

	if len(right.groups) != 1 {
		t.Errorf(`Expected len: %d, received: %d`, 1, len(right.groups))
	}

	entries = left.getEntriesAtValue(left.median())
	checkEntries(t, entries, newCoordinate(0, 0), newCoordinate(0, 1))
	if !left.isLastValue() {
		t.Errorf(`Expected last value.`)
	}

	entries = right.getEntriesAtValue(right.median())
	checkEntries(t, entries, newCoordinate(1, 0), newCoordinate(1, 1))
	if !right.isLastValue() {
		t.Errorf(`Expected last value.`)
	}
}
