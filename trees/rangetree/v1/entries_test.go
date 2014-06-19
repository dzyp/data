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

func checkEntries(t *testing.T, entries []r.Entry, expected ...r.Entry) {
	if len(entries) != len(expected) {
		t.Errorf(`Expected len: %d, received: %d`, len(expected), len(entries))
		return
	}

	for i, entry := range entries {
		if entry.GetDimensionalValue(1) != expected[i].GetDimensionalValue(1) {
			t.Errorf(
				`Expected: %+v, received: %+v`,
				expected[i].GetDimensionalValue(2),
				entry.GetDimensionalValue(1),
			)
		}

		if entry.GetDimensionalValue(2) != expected[i].GetDimensionalValue(2) {
			t.Errorf(
				`Expected: %+v, received: %+v`,
				expected[i].GetDimensionalValue(2), entry.GetDimensionalValue(2),
			)
		}
	}
}

func checkLen(t *testing.T, entries []r.Entry, expected int) {
	if len(entries) != expected {
		t.Errorf(`Expected len: %d, received: %d`, expected, len(entries))
	}
}
