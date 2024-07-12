// Package tools provides common functions.
package tools

func insertElement[T any](array []T, value T, index int) []T {
	return append(array[:index], append([]T{value}, array[index:]...)...)
}

func removeElement[T any](array []T, index int) []T {
	return append(array[:index], array[index+1:]...)
}

func MoveElement[T any](array []T, srcIndex int, dstIndex int) []T {
	value := array[srcIndex]
	return insertElement(removeElement(array, srcIndex), value, dstIndex)
}
