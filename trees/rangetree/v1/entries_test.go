package v1

import (
	"log"
	"reflect"
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

func TestRemove(t *testing.T) {
	log.Printf(`test`)
	p1 := newPoint(0, 0)

	entries := Entries([]r.Entry{p1})

	entries = entries.Remove(p1)

	if len(entries) != 0 {
		t.Errorf(`Expected len: %d, received: %d`, 0, len(entries))
	}

	p2 := newPoint(1, 1)

	entries = append(entries, p1, p2)

	if len(entries) != 2 {
		t.Errorf(`Expected len: %d, received: %d`, 2, len(entries))
	}

	entries = entries.Remove(p1, p2)

	if len(entries) != 0 {
		t.Errorf(`Expected len: %d, received: %d`, 0, len(entries))
	}

	entries = append(entries, p1, p2)
	expected := Entries([]r.Entry{p1})
	entries = entries.Remove(p2)

	if !reflect.DeepEqual(expected, entries) {
		t.Errorf(`Expected: %+v, received: %+v`, expected, entries)
	}
}
