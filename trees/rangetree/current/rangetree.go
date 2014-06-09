package current

import (
	rt "github.com/dzyp/data/trees/rangetree"
	"github.com/dzyp/data/trees/rangetree/v1"
)

func New(maxDimensions int, entries ...rt.Entry) rt.RangeTree {
	return v1.New(maxDimensions, entries...)
}
