package rangetree

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func TestConstructTree(t *testing.T) {
	ints := []int{3, 5, 1, 2, 9, 0, 13}
	New(ints)
}

func TestGetRange(t *testing.T) {
	ints := []int{3, 5, 1, 2, 9, 0, 13}
	tree := New(ints)

	result := tree.GetRange(0, 4)
	if !reflect.DeepEqual(ints[0:4], result) {
		t.Errorf(`Expected result: %+v, received: %+v`, ints[0:5], result)
	}
}

func TestGetRangeMatchingEndpoints(t *testing.T) {
	ints := []int{3, 5, 1, 2, 9, 0, 13}
	tree := New(ints)

	result := tree.GetRange(1, 9)
	if !reflect.DeepEqual(ints[1:5], result) {
		t.Errorf(`Expected result: %+v, received: %+v`, ints[1:5], result)
	}
}

func TestBenchmark(t *testing.T) {
	numInts := 1000000

	ints := make([]int, numInts)
	intMap := make(map[int]bool)
	for i := 0; i < numInts; i++ {
		ints[i] = i
		intMap[i] = true
	}

	tree := New(ints)

	t1 := time.Now()
	result := tree.GetRange(0, numInts)

	log.Printf(
		`It took %d ms to fetch %d items.`,
		time.Since(t1).Nanoseconds()/int64(time.Millisecond),
		numInts,
	)

	log.Println(`Length: %+v`, len(result))

	/*
		if !reflect.DeepEqual(ints, result) {
			t.Errorf(`Expected result: %+v, received: %+v`, ints, result)
		}*/

	t2 := time.Now()
	rangeInts := make([]int, 0)
	for i := 0; i < numInts; i++ {
		if _, ok := intMap[i]; ok {
			rangeInts = append(rangeInts, i)
		}
	}

	log.Printf(
		`It took %d ms to fetch %d items.`,
		time.Since(t2).Nanoseconds()/int64(time.Millisecond),
		numInts,
	)
}
