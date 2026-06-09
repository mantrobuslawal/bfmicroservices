package sliceutil

// Map converts a slice of A values into a slice of B values.
func Map[A any, B any, E error](items []A, mapper func(A) (B, error)) ([]B, error) {
	if len(items) == 0 {
		return nil, nil
	}

	out := make([]B, 0, len(items))

	for _, item := range items {
		mapped, err := mapper(item)
		if err != nil {
			return nil, err
		}
		out = append(out, mapped)
	}

	return out, nil
}
