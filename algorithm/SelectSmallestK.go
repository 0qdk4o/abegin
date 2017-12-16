package main

import (
	"fmt"
)

// Tow array A and B have same count elements as N, elements come from A add elements
// from B produce new array that having N*N count elements, how to select smallest
// K elements from new array
// VIA: google interview? maybe

// suppose:
//    A: 2, 6, 4, 7, 15, 9
//    B: 5, 12, 11, 4, 9, 13
//    K: 10
// so:
//    N = 6
//    len(sum array) = 6*6 = 36
//
// select K elements from A and B and sort these elemens
//    * K >= N, sort A and B directly
//    * K < N, select smallest K elements from A and B, then sort
//
// sum array looks like:
// 2+4   2+5  2+9  2+11  2+12  2+13
// 4+4   4+5  4+9  4+11  4+12  4+13
// 6+4   6+5  6+9  6+11  6+12  6+13
// 7+4   7+5  7+9  7+11  7+12  7+13
// 9+4   9+5  9+9  9+11  9+12  9+13
// 15+4 15+5 15+9 15+11 15+12 15+13
//
// A[0]+B[0] = 2+4, this element is smallest value in the new sum array
// from left to right in the same row, the sum value is increase by degrees
// from top to bottom in the same column, the sum value is increase by degrees
// But it can not say that values from diffrent row and column are increase
// by degrees. such as A[0] + B[2] = 2+9, A[2] + B[1] = 6+5. we should compare
// values from diffrent row and column and find next smallest element from sum
// array
func nextSmallest(a, b []int, c *[]int64) (res int64) {
	csize := len(*c)
	if csize < 1 {
		panic("invalid compare array")
	}

	var sum int
	for n := csize - 1; n >= 0; n-- {
		tempi := (*c)[n] >> 32
		tempj := (*c)[n] & 0xffffffff
		if n == csize-1 {
			sum = a[tempi] + b[tempj]
			continue
		}
		if sum > a[tempi]+b[tempj] {
			sum = a[tempi] + b[tempj]
			(*c)[n], (*c)[csize-1] = (*c)[csize-1], (*c)[n]
		}
	}
	res = (*c)[csize-1]
	*c = (*c)[:csize-1]
	return
}

func appendToComarr(c *[]int64, v int64) {
	for i := 0; i < len(*c); i++ {
		if (*c)[i] == v {
			return
		}
	}
	*c = append(*c, v)
}

func selectSumK(a, b []int, k int) []int {
	if len(a) != len(b) {
		return nil
	}
	size := len(a)
	sumArray := make([]int, k)
	count := 1
	i := 0
	j := 0
	sumArray[0] = a[i] + b[j]
	comarr := make([]int64, 0, k)

	for count < k {
		if i+1 < size {
			appendToComarr(&comarr, int64(i+1)<<32|int64(j))
		}
		if j+1 < size {
			appendToComarr(&comarr, int64(i)<<32|int64(j+1))
		}

		tempIndex := nextSmallest(a, b, &comarr)
		i = int(tempIndex >> 32)
		j = int(tempIndex & 0xffffffff)
		sumArray[count] = a[i] + b[j]

		count++
	}
	return sumArray
}

// quick select algorithm for sort
func qsort(a []int, start, end int) {
	if start >= end {
		return
	}
	if end >= len(a) {
		end = len(a) - 1
	}
	pivot := a[(start+end)/2]

	i := start
	j := end
	for {
		for ; i <= end && a[i] < pivot; i++ {
		}
		for ; j >= start && a[j] > pivot; j-- {
		}
		if i >= j {
			break
		}
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}
	qsort(a, start, i-1)
	qsort(a, j+1, end)
}

// quick select
// k specify smallest elements count
func quickSelectK(a []int, start, end, k int) {
	if start >= end || k-1 < start || k-1 > end {
		return
	}
	pivot := a[(start+end)/2]

	i := start
	j := end
	for {
		for ; i <= end && a[i] < pivot; i++ {
		}
		for ; j >= start && a[j] > pivot; j-- {
		}
		if i >= j {
			break
		}
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}

	if k-1 < i {
		quickSelectK(a, start, i-1, k)
	} else if k-1 > j {
		quickSelectK(a, j+1, end, k)
	}
}

// output: 6 7 8 9 10 11 11 11 12 13
func main() {
	arrA := []int{2, 6, 4, 7, 15, 9}
	arrB := []int{5, 12, 11, 4, 9, 13}
	k := 10
	quickSelectK(arrA, 0, len(arrA)-1, k)
	quickSelectK(arrB, 0, len(arrB)-1, k)
	qsort(arrA, 0, k-1)
	qsort(arrB, 0, k-1)

	sumArr := selectSumK(arrA, arrB, k)
	for i := 0; i < k; i++ {
		fmt.Printf("%d ", sumArr[i])
	}
	fmt.Println()
}
