//
// RMQ(Range Minimum Query)
//

package main

type RMQ struct {
	lcp []int
	st  [][]int
}

func NewRMQ(lcp []int) *RMQ {
	l := len(lcp)
	blen := LOG2[l] + 1

	st, r := make([][]int, blen), make([]int, l)
	for i := 0; i < l; i++ {
		r[i] = i
	}
	st[0] = r

	for p := 1; p < blen; p++ {
		r := make([]int, l)
		for i := 0; i < l; i++ {
			x, y := st[p-1][i], st[p-1][Min(l-1, i+(1<<uint(p-1)))]
			if lcp[x] < lcp[y] {
				r[i] = x
			} else {
				r[i] = y
			}
		}
		st[p] = r
	}

	return &RMQ{lcp, st}
}

func (rmq *RMQ) Query(x, y int) int {
	if y < x {
		x, y = y, x
	}
	if y -= 1; x == y {
		return x
	}

	k := LOG2[y-x]
	i, j := rmq.st[k][x], rmq.st[k][y-(1<<uint(k))+1]

	if rmq.lcp[i] < rmq.lcp[j] {
		return i
	} else {
		return j
	}
}
