// Package sliceutil contains small, generic helpers for working with slices.
//
// The package provides lightweight utilities for common slice transformations
// that are not currently covered by the Go standard library. It is intended for
// simple, explicit operations such as mapping one slice type into another while
// preserving order.
//
// It is the intention of this package to remain small.Pefer plain Go loops for complex
// transformations, including elements such as branching, error handling and similar.
//
// Typical responsibilities include:
//
//   - mapping a slice of one type into a slice of another type;
//   - avoiding repeated allocation bolilerplate in simple adapter code;
//   - keeping transport/domain mapper functions concise where the mapping
//     logic itself is straightforward.
//
// Helpers in this package avoid hidden side effects, and don't mutate input
// slices unless explicitly documented, they also preserve the order of input
// elements.
package sliceutil
