package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "f", "", "Path to the Netscape bookmark file. Reads from stdin if not provided.")
	flag.Parse()

	var reader io.Reader
	var err error

	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", inputFile, err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		// Check if stdin has data
		stat, _ := os.Stdin.Stat()
		if (stat.Mode()&os.ModeCharDevice) != 0 && stat.Size() == 0 {
			fmt.Fprintln(os.Stderr, "Error: No input file provided and no data on stdin.")
			fmt.Fprintln(os.Stderr, "Usage: netscape-bookmarks-json-converter -f <filepath>")
			fmt.Fprintln(os.Stderr, "   or: cat bookmarks.html | netscape-bookmarks-json-converter")
			flag.PrintDefaults()
			os.Exit(1)
		}
		reader = os.Stdin
	}

	bookmarks, err := ParseNetscapeBookmarks(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing bookmarks: %v\n", err)
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}
