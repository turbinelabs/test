package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/turbinelabs/test/server/http"
)

func main() {
	ts, err := http.NewTestServerFromFlagSet(flag.CommandLine, os.Args[1:])
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Usage: ")
		flag.CommandLine.PrintDefaults()
		os.Exit(1)
	}
	ts.Serve()
}
