package rangetree

type Entry interface {
	/*
		Pass in an int representing the dimension of interest and returns
		a value to be sorted on in that dimension.
	*/
	GetDimensionalValue(dimension int) int
	/*
		The number of dimensions held by this entry.
	*/
	MaxDimensions() int
	/*
		Returns a bool indicating whether values are equal up to the given dimension
	*/
	EqualAtDimension(entry Entry, dimension int) bool
	/*
		Returns a value indicating relationship at the given dimension
	*/
	LessThan(entry Entry, dimension int) bool
}

type Bounds interface {
	/*
		[Low, High) Houses the high/low values for a query
	*/
	High() int
	Low() int
}

type Query interface {
	/*
		Returns a bounds interface for the given dimension.  Return
		nil if the dimension is outside the bounds of the query
	*/
	GetDimensionalBounds(dimension int) Bounds
}

type RangeTree interface {
	Remove(entries ...Entry)
	GetRange(query Query) []Entry
	Insert(entries ...Entry)
	Copy() RangeTree
	Clear()
}
