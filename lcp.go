//
// LCP(Longest Common Prefix)
//

package main

import (
	"sync"

	"github.com/tobi-c/go-sais/suffixarray"
)

type RepeatHistogram map[string][]int64
type RepeatHistogramMap struct {
	sync.RWMutex
	RepHist RepeatHistogram
}

type Result struct {
	From, Length, Count int
}

func LCPDevide(read string) {
	rread, l := Reverse(read), len(read)+1
	lcp, rank := LCPArray(read)
	rlcp, rrank := LCPArray(rread)
	rmq, rrmq := NewRMQ(lcp), NewRMQ(rlcp)

	UpsertRepeat(read,
		LCPDevideRecursive(l, lcp, rank, rmq, rlcp, rrank, rrmq, 0, l))

	return
}

func UpsertRepeat(read string, r []Result) {
	if len(r) <= 1 {
		return
	}

	// Remove duplicates
	rm := make(map[Result]bool)
	for _, v := range r {
		rm[v] = true
	}

	// Remove repeats
	for v := range rm {
		if IsRepeatString(read[v.From:(v.From + v.Length)]) {
			delete(rm, v)
		}
	}

	var apps []ApproxRepeat
	if Approx {
		apps = ExtendRepeat(read, rm)
	}

	Repeats.Lock()
	for v := range rm {
		r := read[v.From:(v.From + v.Length)]
		l := v.Length * v.Count
		if v, ok := Repeats.RepHist[r]; ok {
			v[l]++
		} else {
			Repeats.RepHist[r] = make([]int64, HistogramWidth)
			Repeats.RepHist[r][l] = 1
		}
	}

	if Approx && len(apps) > 0 {
		for _, app := range apps {
			r, l := app.Repeat, app.Length
			if v, ok := Repeats.RepHist[r]; ok {
				v[l]++
			} else {
				Repeats.RepHist[r] = make([]int64, HistogramWidth)
				Repeats.RepHist[r][l] = 1
			}
		}
	}
	Repeats.Unlock()
	rm = nil

	return
}

func LCPDevideRecursive(
	l int, lcp, rank []int, rmq *RMQ,
	rlcp, rrank []int, rrmq *RMQ, start, end int) []Result {

	if end-start <= 1 {
		return []Result{}
	}

	r := make([]Result, 0, l/2)
	p := (start + end) / 2
	for j := 1; j < (end - p + 2); j++ {
		if i := p + j; i < end {
			if lcpxy := lcp[rmq.Query(rank[p], rank[i])]; lcpxy > 0 {
				if i2 := p + lcpxy + j - 1; i2 < end {
					if i2 = l - 2 - i2; i2 >= 0 {
						if cnt := rlcp[rrmq.Query(rrank[i2], rrank[i2+j])] / j; cnt > 0 {
							cnt++
							r = append(r, Result{l - 2 - (i2 + cnt*j - 1), j, cnt})
						}
					}
				}
			}
		}

		if i := l - 2 - p; i >= 0 && p >= start+j {
			if lcpxy := rlcp[rrmq.Query(rrank[i], rrank[i+j])]; lcpxy > 0 {
				if e := l - i - j - lcpxy - 1; e >= start {
					if cnt := lcp[rmq.Query(rank[e], rank[e+j])] / j; cnt > 0 {
						r = append(r, Result{e, j, cnt + 1})
					}
				}
			}
		}
	}

	left := LCPDevideRecursive(
		l, lcp, rank, rmq, rlcp, rrank, rrmq, start, p)

	right := LCPDevideRecursive(
		l, lcp, rank, rmq, rlcp, rrank, rrmq, p+1, end)

	return append(append(left, r...), right...)
}

func LCPArray(read string) ([]int, []int) {
	read += "\x00"
	sarr := suffixarray.New([]byte(read)).SuffixArray()
	size := len(sarr)
	rank := make([]int, size)
	for i := 0; i < size; i++ {
		rank[sarr[i]] = i
	}

	larr := make([]int, size)
	for i, lcp := 0, 0; i < size; i++ {
		idx := rank[i]
		pos1 := sarr[idx]
		if idx == size-1 {
			lcp = 0
			continue
		}
		pos2 := sarr[idx+1]
		for i, j := pos1+lcp, pos2+lcp; read[i] == read[j]; i, j, lcp = i+1, j+1, lcp+1 {
		}
		larr[idx] = lcp

		if lcp--; lcp < 0 {
			lcp = 0
		}
	}

	return larr, rank
}
