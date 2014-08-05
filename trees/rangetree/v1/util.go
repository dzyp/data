package v1

func splitNodes(slice []*node, numParts int) [][]*node {
	parts := make([][]*node, numParts)
	for i := 0; i < numParts; i++ {
		parts[i] = slice[i*len(slice)/numParts : (i+1)*len(slice)/numParts]
	}
	return parts
}
