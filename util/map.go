package util

func Map[I any, R any](input []I, mapper func(I, int) R) []R {
	results := make([]R, len(input))

	for idx, elem := range input {
		results[idx] = mapper(elem, idx)
	}

	return results
}
