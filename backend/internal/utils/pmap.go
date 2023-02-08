package utils

import "sync"

func ParallelMap[T1 any, T2 any](arr []T1, fn func(int, T1) T2) []T2 {
	wg := &sync.WaitGroup{}
	wg.Add(len(arr))

	output := make([]T2, len(arr), len(arr))

	for i := range arr {
		go func(index int, x T1) {
			defer wg.Done()

			result := fn(index, x)
			output[index] = result

		}(i, arr[i])
	}

	wg.Wait()
	return output
}
