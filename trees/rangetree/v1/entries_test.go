package v1

import (
	"testing"

	r "github.com/dzyp/data/trees/rangetree"
)

func checkSortedList(t *testing.T, entries []int, expected ...int) {
	if len(entries) != len(expected) {
		t.Errorf(`Expected len: %d, received: %d`, len(expected), len(entries))
		return // don't want to panic
	}

	for i, entry := range entries {
		coord := expected[i]

		if entry != coord {
			t.Errorf(`Expected x: %d, received: %d`, coord, entry)
		}

		if entry != coord {
			t.Errorf(`Expected y: %d, received: %d`, coord, entry)
		}
	}
}

func TestEntriesSortOnFirstDimension(t *testing.T) {
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

func TestEntriesSortOnSecondDimension(t *testing.T) {
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

	checkSortedList(t, sorted, 0, 1)

	/*
		entries = ew.getEntriesAtValue(sorted[0])
		checkEntries(t, entries, newCoordinate(0, 0), newCoordinate(0, 1))

		entries = ew.getEntriesAtValue(sorted[1])
		checkEntries(t, entries, newCoordinate(1, 0), newCoordinate(1, 1))*/
}

func TestEntriesWrapperOneValue(t *testing.T) {
	entries := []r.Entry{
		newPoint(0, 1),
	}

	ew := newEntries(entries, 1, true)

	sorted := ew.getSortedValues()
	checkSortedList(t, sorted, 0)

	entries = ew.getEntriesAtValue(sorted[0])
	checkEntries(t, entries, newCoordinate(0, 1))
}

func TestEntriesWrapperSplit(t *testing.T) {
	entries := []r.Entry{
		newPoint(0, 1),
		newPoint(1, 1),
		newPoint(0, 0),
		newPoint(1, 0),
	}

	ew := newEntries(entries, 1, true)
	entry := ew.median()

	checkSortedList(t, []int{entry}, 1)

	left, right := ew.split(-1)

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
