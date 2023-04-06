package utils

import "sync"

func Map[T1 any, T2 any](arr *[]T1, f func(T1) T2) *[]T2 {
	result := make([]T2, len(*arr))
	for i, elem := range *arr {
		result[i] = f(elem)
	}
	return &result
}

func AsyncMap[T1 any, T2 any](arr *[]T1, f func(T1) T2) *[]T2 {
	result := make([]T2, len(*arr))
	var wg sync.WaitGroup
	wg.Add(len(*arr))
	for i, elem := range *arr {
		go func(i int, elem T1) {
			result[i] = f(elem)
			wg.Done()
		}(i, elem)
	}
	wg.Wait()
	return &result
}

func Filter[T any](arr []T, fun func(T) bool) []T {
	dst := make([]T, 0, len(arr))
	for _, v := range arr {
		if fun(v) {
			dst = append(dst, v)
		}
	}
	return dst
}
