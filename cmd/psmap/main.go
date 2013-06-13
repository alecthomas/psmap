package main

import (
	"encoding/csv"
	"fmt"
	"github.com/alecthomas/psmap"
	flag "github.com/ogier/pflag"
	"io"
	"os"
)

var (
	csvFlag = flag.Bool("csv", false, "build psmap from CSV file")
)

func Fatalf(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+f+"\n", args...)
	os.Exit(1)
}

func main() {
	flag.Usage = func() {
		fmt.Print(`psmap: Build/dump psmaps.

Build:
    psmap --csv <csv> <psmap>

Dump:
    psmap <psmap>

Flags:
`)
		flag.PrintDefaults()
	}
	flag.Parse()
	if *csvFlag {
		if flag.NArg() != 2 {
			Fatalf("need input and output filenames")
		}
		inf := flag.Arg(0)
		outf := flag.Arg(1)
		build(inf, outf)
	} else {
		if flag.NArg() != 1 {
			Fatalf("need input filename")
		}
		inf := flag.Arg(0)
		dump(inf)
	}
}

func build(inf, outf string) {
	r, err := os.Open(inf)
	if err != nil {
		Fatalf("could not open %s: %s", inf, err)
	}
	defer r.Close()
	w, err := os.Create(outf)
	if err != nil {
		Fatalf("could not create %s: %s", outf, err)
	}
	defer w.Close()
	writer := psmap.NewBuilder(w)
	defer writer.Close()
	reader := csv.NewReader(r)
	n := 1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			Fatalf("%s: %d: failed to read CSV file: %s", inf, n, err)
		}
		if len(row) < 2 {
			Fatalf("%s: %d: CSV files must have at least 2 columns, found %d", inf, n, len(row))
		}
		writer.Add([]byte(row[0]), []byte(row[1]))
		n++
	}
}

func dump(inf string) {
	pm, err := psmap.Open(inf)
	if err != nil {
		Fatalf("failed to open psmap %s: %s", inf, err)
	}
	for kv := range pm.Iterate() {
		fmt.Printf("%s %s\n", kv.Key, kv.Value)
	}
}
