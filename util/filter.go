package util

func Filter[I any](input []I, filter func(I, int) bool) []I {
	results := []I{}

	for idx, elem := range input {
		if filter(elem, idx) {
			results = append(results, elem)
		}
	}

	return results
}
