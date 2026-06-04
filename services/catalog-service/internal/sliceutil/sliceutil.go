package sliceutil

// Map converts a slice of A values into a slice of B values.
func Map[A any, B any] (items []A, mapper func(A) B) []B {
	if len(items) == 0 {
		return nil
	}

	out := make([]B, 0, len(items))

	for _, item := range items {
		out = append(out, mapper(item))
	}

	return out
}

