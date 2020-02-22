package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	gzip "github.com/klauspost/pgzip"
)

const (
	PARALLEL                 int = 8
	LOG2LIST_LENGTH          int = 1000
	FASTQ_FILENAME_EXTENSION     = `_?R[12].*fastq(\.gz)?$`
	HISTOGRAM_EXACT_MATCH        = "_histogram_exact.lst"
	HISTOGRAM_APPROX_MATCH       = "_histogram_approx.lst"
)

var (
	Approx         bool
	Repeats        *RepeatHistogramMap
	HistogramWidth int = 0
	ExitCode       int = 0
	POW2           []int
	LOG2           []int
)

func init() {
	LOG2 = make([]int, LOG2LIST_LENGTH)
	for i := 1; i < LOG2LIST_LENGTH; i++ {
		LOG2[i] = int(math.Log2(float64(i)))
	}

	// Don't print timestamp
	log.SetFlags(0)
}

func main() {
	defer func() { os.Exit(ExitCode) }()

	// Parse option switches
	fastq1, fastq2, output, threads := parseCommandLine()

	f1, r1 := SetupFile(fastq1)
	defer f1.Close()
	defer r1.Close()

	rs := make([]*gzip.Reader, 0, 2)
	rs = append(rs, r1)
	f2, r2 := &os.File{}, &gzip.Reader{}
	if len(fastq2) > 0 {
		f2, r2 = SetupFile(fastq2)
		rs = append(rs, r2)
		defer f2.Close()
		defer r2.Close()
	}

	Repeats = new(RepeatHistogramMap)
	Repeats.RepHist = make(RepeatHistogram)
	CountRepeats(threads, rs)
	OutputHistgramData(output)
}

func CountRepeats(threads int, rs []*gzip.Reader) {
	var wg1 sync.WaitGroup
	wg1.Add(len(rs))
	for _, r := range rs {
		go func(ir io.Reader) {
			defer wg1.Done()
			ch := make(chan string, PARALLEL*100)
			var wg2 sync.WaitGroup
			wg2.Add(threads)
			for i := 0; i < threads; i++ {
				go func(read <-chan string) {
					defer wg2.Done()
					for rd := range read {
						LCPDevide(rd)
					}
				}(ch)
			}

			sc := bufio.NewScanner(ir)

			// calculate histogram width by first read length
			if HistogramWidth == 0 {
				sc.Scan()
				sc.Scan()
				HistogramWidth = len(sc.Text()) * 2
				ch <- sc.Text()
				sc.Scan()
				sc.Scan()
			}

			// continue scanning
			for sc.Scan() {
				sc.Scan()
				ch <- sc.Text()
				sc.Scan()
				sc.Scan()
			}
			close(ch)
			wg2.Wait()
			checkError("Read a fastq file", sc.Err())
		}(r)
	}
	wg1.Wait()
}

func OutputHistgramData(output string) {
	out, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	checkError("Open a output file", err)
	defer out.Close()

	for k, v := range Repeats.RepHist {
		fmt.Fprintln(out, CreateCsvLine(k, v))
	}
}

func parseCommandLine() (string, string, string, int) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:

  %s [options] fastq1 [fastq2]

`, filepath.Base(os.Args[0]), os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	var (
		output  string
		threads int
	)
	flag.BoolVar(&Approx, "a", false, "Approximate matching")
	flag.BoolVar(&Approx, "approximateMatching", false, "Approximate matching(long option)")
	flag.IntVar(&threads, "P", PARALLEL, "Number of threads")
	flag.IntVar(&threads, "numberOfThreads", PARALLEL, "Number of threads(long option)")
	flag.IntVar(&HistogramWidth, "L", 0, "Histogram width")
	flag.IntVar(&HistogramWidth, "histogramWidth", 0, "Histogram width(long option)")
	flag.StringVar(&output, "o", "", "Output file")
	flag.StringVar(&output, "outputFile", "", "Output file(long option)")
	flag.Parse()

	f := flag.Args()
	l := len(f)
	if l < 1 {
		log.Print("No fastq files.")
		flag.Usage()
		os.Exit(ExitCode)
	}

	fastq1 := f[0]
	if _, err := os.Stat(fastq1); os.IsNotExist(err) {
		checkError("Fastq file", err)
	}

	fastq2 := ""
	if l > 1 {
		fastq2 = f[1]
		if _, err := os.Stat(fastq2); os.IsNotExist(err) {
			checkError("Fastq file", err)
		}
	}

	if threads < 1 {
		threads = PARALLEL
	}

	if HistogramWidth > 0 {
		HistogramWidth++
	}

	if output == "" {
		cwd, err := os.Getwd()
		checkError("Current Directory", err)
		base := filepath.Base(fastq1)
		output = filepath.Join(cwd,
			regexp.MustCompile(FASTQ_FILENAME_EXTENSION).ReplaceAllString(base, ""))
		if Approx {
			output += HISTOGRAM_APPROX_MATCH
		} else {
			output += HISTOGRAM_EXACT_MATCH
		}
	}

	return fastq1, fastq2, output, threads
}
