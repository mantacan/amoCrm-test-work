package main

import "fmt"

// PrintSorted принимает на вход два канала, каждый из которых возвращает конечную монотонно неубывающую
// последовательность целых чисел (т.е. отсортированные по возрастанию). Необходимо объединить значения
// из обоих каналов и вывести их в stdout в виде одной монотонно неубывающей последовательности.
//
// Пример:
// In: a = [0, 0, 3, 4]; b = [1, 3, 4, 6, 8]
// Out: res = [0, 0, 1, 3, 3, 4, 4, 6, 8]
//
// Для проверки решения запустите тесты: go test -v

func ComputeSortedArray(channels ...<-chan int) ([]int, error) {
	type val struct {
		v    int
		prev int
	}
	channelData := make(map[<-chan int]val)
	result := make([]int, 0)

	for _, ch := range channels {
		if v, ok := <-ch; ok {
			channelData[ch] = val{v: v, prev: v}
		}
	}

	for len(channelData) > 0 {
		var minCh <-chan int

		for ch := range channelData {
			if minCh == nil || channelData[ch].v < channelData[minCh].v {
				minCh = ch
			}
		}

		result = append(result, channelData[minCh].v)

		if v, ok := <-minCh; ok {
			if v < channelData[minCh].prev {
				return nil, fmt.Errorf("channel %d: unsorted data: after %d", v, channelData[minCh].prev)
			}
			data := channelData[minCh]
			channelData[minCh] = val{v, data.v}
		} else {
			delete(channelData, minCh)
		}
	}

	return result, nil
}

func PrintSorted(channels ...<-chan int) {
	arr, err := ComputeSortedArray(channels...)
	if err == nil {
		for _, v := range arr {
			fmt.Println(v)
		}
	} else {
		fmt.Println(err)
	}
}

func main() {
	a := make(chan int)
	b := make(chan int)

	go func() {
		defer close(a)
		a <- 1
		a <- 4
		a <- 6
	}()

	go func() {
		defer close(b)
		b <- 2
		b <- 3
		b <- 5
		b <- 7
		b <- 8
		b <- 9
	}()

	PrintSorted(a, b)
}
