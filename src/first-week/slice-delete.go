package main

import "fmt"

func deleteAt[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice
	}
	copy(slice[index:], slice[index+1:])
	newSlice := slice[:len(slice)-1]
	if cap(newSlice) > 2*len(newSlice) {
		newSlice = append([]T(nil), newSlice...) // 缩容
	}
	return newSlice
}

func main() {
	slice := []int{1, 2, 3, 4, 5}
	fmt.Println("Before:", slice, "cap:", cap(slice))
	slice = deleteAt(slice, 2)
	fmt.Println("After:", slice, "cap:", cap(slice))

	strSlice := []string{"a", "b", "c", "d", "e"}
	fmt.Println("Before:", strSlice, "cap:", cap(strSlice))
	strSlice = deleteAt(strSlice, 2)
	fmt.Println("After:", strSlice, "cap:", cap(strSlice))
}
