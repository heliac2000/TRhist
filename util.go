//
// Utility functions
//

package main

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"strconv"

	gzip "github.com/klauspost/pgzip"
)

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Min(values ...int) int {
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}

	return min
}

func IsRepeatString(r string) bool {
	l := len(r)
	if l <= 1 {
		return false
	}

L:
	for i := (l - 1); i >= 1; i-- {
		if l%i != 0 {
			continue
		}
		x := SplitByNbytes(r, i)
		for j := 1; j < len(x); j++ {
			if x[j] != x[0] {
				continue L
			}
		}
		return true
	}

	return false
}

func checkError(msg string, err error) {
	if err != nil {
		log.Print(msg+": ", err)
		ExitCode = 1
		runtime.Goexit()
	}
}

func SetupFile(fq string) (*os.File, *gzip.Reader) {
	f, err := os.Open(fq)
	checkError("Open a fastq file", err)
	r, err := gzip.NewReader(bufio.NewReader(f))
	checkError("Gzip compressed stream", err)

	return f, r
}

func CreateCsvLine(key string, value []int64) string {
	line := key
	for _, v := range value[2:] {
		line += "," + strconv.FormatInt(v, 10)
	}

	return line
}

// n must be equal or greater than 0
func Pow2(n int) int {
	return POW2[n]
}

// n must be greater than 0
func Log2(n int) int {
	return LOG2[n]
}

func Ascend(x, y int) (int, int) {
	if x > y {
		return y, x
	}
	return x, y
}

func SplitByNbytes(str string, w int) []string {
	subs := make([]string, 0, len(str)/w)
	var sub string
	for len(str) >= w {
		sub, str = str[:w], str[w:]
		subs = append(subs, sub)
	}

	return subs
}
