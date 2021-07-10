package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var parallel int

	flag.Usage = func() { usage() }
	flag.IntVar(&parallel, "parallel", 10, "Number of requests to send in parallel")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage()
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
