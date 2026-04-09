package main

import (
	"os"

	"github.com/teoc98/shp/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		// error is discarded as cobra already reported it
		os.Exit(1)
	}
}
