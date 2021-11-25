package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"

	"github.com/shipyardapp/gooptions/model"
)

var output = flag.String("output", "", "The output file name, or empty to use stdout.")

func main() {
	flag.Parse()

	structType, err := model.NewStructTypeFromReflectType(ReflectTypeVar)
	if err != nil {
		exit(err, 1)
	}

	outputFile := os.Stdout
	if *output != "" {
		var err error
		outputFile, err = os.Create(*output)
		if err != nil {
			exit(err, 2)
		}
		defer func() {
			if err := outputFile.Close(); err != nil {
				exit(fmt.Errorf("failed to close output file: %v", err), 3)
			}
		}()
	}

	if err := gob.NewEncoder(outputFile).Encode(structType); err != nil {
		exit(fmt.Errorf("gob encode error: %v", err), 4)
	}
}

func exit(err error, exitCode int) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(exitCode)
}
