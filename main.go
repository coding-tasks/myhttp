package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var parallel int

	flag.Usage = func() { usage() }
	flag.IntVar(&parallel, "parallel", defaultParallelRequests, "Number of requests to send in parallel")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage()
	}

	dl := NewDownloader(args, WithParallelRequests(parallel))
	for hash := range dl.Download() {
		if hash.err != nil {
			fmt.Printf("%s %s\n", hash.url, hash.err)
		} else {
			fmt.Printf("%s %s\n", hash.url, hash.sum)
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `
USAGE
  myhttp URL<1...N> [FLAGS]

FLAGS
  -parallel int	Number of requests to send in parallel, defaults to 10

ARGUMENTS
  URL		Website url, separated by space, to fetch contents from

EXAMPLES
  $ myhttp google.com
  $ myhttp -parallel 3 google.com facebook.com yahoo.com

`)
	os.Exit(1)
}
