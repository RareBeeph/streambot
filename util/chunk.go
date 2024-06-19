package util

// Chunk splits a slice into smaller chunks of the specified size.
// It takes a slice of items and a chunkSize as input parameters and returns
// a slice of slices, where each inner slice is a chunk.
//
// When len(items)%chunkSize!=0, the length of the final chunk is instead
// the remainder.
//
// Credit to https://stackoverflow.com/a/72408490
//
// Example:
//
//	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
//	chunks := Chunk(items, 3)
//	// chunks: [[1 2 3] [4 5 6] [7 8 9] [10 11]]
//
// Parameters:
//
//	items: The slice of items to be chunked.
//	chunkSize: The size of each chunk.
//
// Returns:
//
//	chunks: A slice of slices, where each inner slice is a chunk of items.
func Chunk[T any](items []T, chunkSize int) (chunks [][]T) {
	// Don't add an empty inner element
	if len(items) == 0 {
		return
	}

	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
