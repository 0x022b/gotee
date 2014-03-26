package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	appVersion = "0.1.0"
)

func appName() string {
	binary := filepath.Base(os.Args[0])
	return binary[:len(binary)-len(filepath.Ext(binary))]
}

func main() {
	var flags struct {
		append, version *bool
	}
	flags.append = flag.Bool("a", false, "append to the given FILEs, do not overwrite")
	flags.version = flag.Bool("v", false, "output version information and exit")
	flag.Usage = usage
	flag.Parse()

	if *flags.version {
		version()
	}

	mode := os.O_CREATE
	if *flags.append {
		mode |= os.O_APPEND
	}

	writers := []io.Writer{os.Stdout}
	for _, arg := range flag.Args() {
		if file, err := os.OpenFile(arg, mode, 0660); err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			defer file.Close()
			writers = append(writers, file)
		}
	}

	b := make([]byte, 128)
	r := io.Reader(os.Stdin)
	w := io.MultiWriter(writers...)
	for {
		if n, err := r.Read(b); err != nil {
			break
		} else if n > 0 {
			w.Write(b[:n])
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... [FILE]...\n"+
		"Copy standard input to each FILE, and also to standard output.\n\n", appName())
	flag.PrintDefaults()
}

func version() {
	fmt.Fprintf(os.Stderr, "%s %s\n", appName(), appVersion)
	os.Exit(0)
}
